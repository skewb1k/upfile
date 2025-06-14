package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/service"
	"github.com/spf13/cobra"
)

func statusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status [<dir>]",
		Short: "Print status of files in dir (default: current dir)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) == 1 {
				dir = args[0]
			}

			absDir, err := filepath.Abs(dir)
			if err != nil {
				return fmt.Errorf("failed to get abs path to dir: %w", err)
			}

			if err := service.Status(
				cmd.Context(),
				cmd.OutOrStdout(),
				getIndexFsProvider(),
				absDir,
			); err != nil {
				return fmt.Errorf("cannot show status of '%s': %w", absDir, err)
			}
			return nil
		},
	}

	return cmd
}
