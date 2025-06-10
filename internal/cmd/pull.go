package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/entries"
	"github.com/skewb1k/upfile/internal/upstreams"
	"github.com/spf13/cobra"
)

func pullCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull <path>",
		Short: "Pull file from upstream",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("failed to get abs path to dest dir: %w", err)
			}

			fname, destDir := filepath.Base(path), filepath.Dir(path)

			baseDir := getBaseDir()
			upstreamsProvider := upstreams.NewProvider(baseDir)
			entriesProvider := entries.NewProvider(baseDir)

			upstream, err := upstreamsProvider.GetUpstream(cmd.Context(), fname)
			if err != nil {
				return fmt.Errorf("get upstream: %w", err)
			}

			existing, err := os.ReadFile(path)
			if err != nil {
				if !errors.Is(err, os.ErrNotExist) {
					return err
				}
			} else {
				if upstream.Hash.EqualBytes(existing) {
					return ErrUpToDate
				}
			}

			if err := entriesProvider.CreateEntry(cmd.Context(), fname, destDir); err != nil {
				if !errors.Is(err, entries.ErrExists) {
					return fmt.Errorf("create entry: %w", err)
				}
			}

			if err := WriteFile(path, upstream.Content); err != nil {
				return fmt.Errorf("write file: %w", err)
			}

			return nil
		},
	}

	return cmd
}

func WriteFile(path string, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create parent dirs: %w", err)
	}

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
