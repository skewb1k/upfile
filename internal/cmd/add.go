package cmd

import (
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
		RunE: withService(func(cmd *cobra.Command, args []string, s *service.Service) error {
			path, err := filepath.Abs(filepath.Clean(args[0]))
			if err != nil {
				return fmt.Errorf("get abs path to file: %w", err)
			}

			if err := s.Add(cmd.Context(), path); err != nil {
				return err
			}

			cmd.Printf("Added: %s\n", path)

			return nil
		}),
	}
}
