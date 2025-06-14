package cmd

import (
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
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("failed to get abs path to file: %w", err)
			}

			if err := service.Push(
				cmd.Context(),
				cmd.OutOrStdout(),
				getStore(),
				path,
			); err != nil {
				return fmt.Errorf("cannot push '%s': %w", path, err)
			}

			return nil
		},
	}
}
