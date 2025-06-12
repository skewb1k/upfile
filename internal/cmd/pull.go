package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/service"
	"github.com/skewb1k/upfile/internal/store"
	"github.com/spf13/cobra"
)

func pullCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull <path>",
		Short: "Pull file from upstream",
		Args:  cobra.ExactArgs(1),
		RunE: wrap(func(cmd *cobra.Command, s *service.Service, args []string) error {
			path, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("failed to get abs path to dest dir: %w", err)
			}

			fname, destDir := filepath.Base(path), filepath.Dir(path)

			upstream, err := s.CatLatest(cmd.Context(), fname)
			if err != nil {
				return fmt.Errorf("cat latest: %w", err)
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

			if err := s.CreateEntry(cmd.Context(), fname, destDir); err != nil {
				if !errors.Is(err, store.ErrExists) {
					return fmt.Errorf("create entry: %w", err)
				}
			}

			if err := WriteFile(path, upstream); err != nil {
				return fmt.Errorf("write file: %w", err)
			}

			return nil
		}),
	}

	return cmd
}

func WriteFile(path string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create parent dirs: %w", err)
	}

	if err := os.WriteFile(path, content, 0o600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
