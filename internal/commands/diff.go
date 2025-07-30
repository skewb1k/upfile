package commands

import (
	"fmt"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/service"
	"github.com/spf13/cobra"
)

func Diff() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diff <path>",
		Short: "Compare file with its upstream",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("failed to get absolute path to file '%s': %w", args[0], err)
			}

			if err := service.Diff(
				cmd.Context(),
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				cmd.ErrOrStderr(),
				getIndexFsProvider(),
				path,
			); err != nil {
				return fmt.Errorf("cannot show diff with '%s': %w", path, err)
			}

			return nil
		},
	}

	return cmd
}
