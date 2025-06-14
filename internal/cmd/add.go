package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/service"
	"github.com/spf13/cobra"
)

func addCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <path>",
		Short: "Add a file to be tracked",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("failed to get abs path to file: %w", err)
			}

			if err := service.Add(cmd.Context(), getStore(), path); err != nil {
				return fmt.Errorf("cannot add '%s': %w", path, err)
			}

			return nil
		},
	}

	return cmd
}
