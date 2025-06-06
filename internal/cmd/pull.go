package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/service"

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

			if err := s.Pull(cmd.Context(), path); err != nil {
				if errors.Is(err, service.ErrUpToDate) {
					cmd.Println("File up-to-date")

					return nil
				}

				return err //nolint: wrapcheck
			}

			return nil
		}),
	}

	return cmd
}
