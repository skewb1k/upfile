package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/store"
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

			destDir, fname := filepath.Dir(path), filepath.Base(path)

			s := store.New(getBaseDir())

			upstream, err := s.GetUpstream(cmd.Context(), fname)
			if err != nil {
				if errors.Is(err, store.ErrNotFound) {
					return ErrNotTracked
				}

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

			if err := s.CreateEntry(cmd.Context(), fname, destDir); err != nil {
				if !errors.Is(err, store.ErrExists) {
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
