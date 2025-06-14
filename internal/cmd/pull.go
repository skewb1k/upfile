package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/skewb1k/upfile/internal/service"
	"github.com/spf13/cobra"
)

func pullCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:               "pull <filename>",
		Short:             "Pull file from upstream",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completeFname,
		RunE: func(cmd *cobra.Command, args []string) error {
			currentDir, err := os.Getwd()
			if err != nil {
				return err
			}

			if err := service.Pull(
				cmd.Context(),
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				getStore(),
				yes,
				currentDir,
				args[0],
			); err != nil {
				if errors.Is(err, service.ErrCancelled) {
					os.Exit(1)
				}

				return fmt.Errorf("cannot pull '%s': %w", args[0], err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Automatic 'yes' to prompts")

	return cmd
}
