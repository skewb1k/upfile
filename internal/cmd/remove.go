package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/service"

	"github.com/spf13/cobra"
)

func remove() *cobra.Command {
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

			if err := s.Remove(cmd.Context(), path); err != nil {
				return err //nolint: wrapcheck
			}

			cmd.Printf("Removed: %s\n", path)

			return nil
		}),
	}
}
