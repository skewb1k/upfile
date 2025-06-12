package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/service"
	"github.com/spf13/cobra"
)

func removeCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "remove <path>",
		Short:   "Remove entry from tracked list",
		Aliases: []string{"rm"},
		Args:    cobra.ExactArgs(1),
		RunE: wrap(func(cmd *cobra.Command, s *service.Service, args []string) error {
			path, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("failed to get abs path to file: %w", err)
			}

			entryDir, fname := filepath.Dir(path), filepath.Base(path)

			if err := s.DeleteEntry(cmd.Context(), fname, entryDir); err != nil {
				if errors.Is(err, service.ErrNotFound) {
					return ErrNotTracked
				}

				return fmt.Errorf("delete entry: %w", err)
			}

			return nil
		}),
	}
}
