package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"upfile/internal/service"
)

func add() *cobra.Command {
	return &cobra.Command{
		Use:   "add <path>",
		Short: "Add a file to be tracked",
		Args:  cobra.ExactArgs(1),
		RunE: withService(func(cmd *cobra.Command, s *service.Service, args []string) error {
			path, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("failed to get abs path to file: %w", err)
			}

			if err := s.Add(cmd.Context(), path); err != nil {
				if errors.Is(err, service.ErrAlreadyTracked) {
					cmd.PrintErrf("error: file %q already tracked\n", path)
					return nil
				}

				if errors.Is(err, service.ErrFileNotFound) {
					cmd.PrintErrf("error: file %q not found\n", path)
					return nil
				}

				return err //nolint: wrapcheck
			}

			cmd.Printf("Added: %s\n", path)

			return nil
		}),
	}
}
