package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/service"

	"github.com/spf13/cobra"
)

func pushCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "push <path>",
		Short: "Push file to the upstream",
		Args:  cobra.ExactArgs(1),
		RunE: wrap(func(cmd *cobra.Command, s *service.Service, args []string) error {
			path, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("failed to get abs path to file: %w", err)
			}

			if err := s.Push(cmd.Context(), path); err != nil {
				if errors.Is(err, service.ErrUpToDate) {
					cmd.Println("File up-to-date")
					return nil
				}

				return err //nolint: wrapcheck
			}

			return nil
		}),
	}
}
