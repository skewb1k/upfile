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
	var track bool

	cmd := &cobra.Command{
		Use:               "pull <filename>...",
		Short:             "Pull file(s) from upstream",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: completeFnames,
		RunE: func(cmd *cobra.Command, args []string) error {
			currentDir, err := os.Getwd()
			if err != nil {
				return err
			}

			for _, name := range args {
				if err := service.Pull(
					cmd.Context(),
					cmd.InOrStdin(),
					cmd.OutOrStdout(),
					getIndexFsProvider(),
					yes,
					currentDir,
					name,
					track,
				); err != nil {
					if errors.Is(err, service.ErrCancelled) {
						os.Exit(1)
					}

					return fmt.Errorf("cannot pull '%s': %w", name, err)
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Automatic 'yes' to prompts")
	cmd.Flags().BoolVarP(&track, "track", "t", false, "Start tracking pulled files")

	return cmd
}
