package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/store"
	"github.com/spf13/cobra"
)

func removeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove <path>",
		Short:   "Remove entry from tracked list",
		Aliases: []string{"rm"},
		Args:    cobra.ExactArgs(1),
		RunE: withStore(func(cmd *cobra.Command, s *store.Store, args []string) error {
			path, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("failed to get abs path to file: %w", err)
			}

			entry, fname := filepath.Dir(path), filepath.Base(path)

			if err := s.DeleteEntry(cmd.Context(), fname, entry); err != nil {
				if errors.Is(err, store.ErrNotFound) {
					return ErrNotTracked
				}

				return fmt.Errorf("delete entry: %w", err)
			}

			return nil
		}),
	}

	return cmd
}
