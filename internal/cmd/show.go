package cmd

import (
	"fmt"

	"github.com/skewb1k/upfile/internal/service"
	"github.com/spf13/cobra"
)

func showCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "show <filename>",
		Short:             "Show upstream version of file",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completeFname,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := service.Show(
				cmd.Context(),
				cmd.OutOrStdout(),
				getStore(),
				args[0],
			); err != nil {
				return fmt.Errorf("cannot show '%s': %w", args[0], err)
			}

			return nil
		},
	}

	return cmd
}
