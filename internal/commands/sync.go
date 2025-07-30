package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/skewb1k/upfile/internal/service"
	"github.com/spf13/cobra"
)

func Sync() *cobra.Command {
	// TODO:
	// like git add -p
	// var patch bool

	var yes bool

	cmd := &cobra.Command{
		Use:               "sync <filename>",
		Args:              cobra.ExactArgs(1),
		Short:             "Sync all entries of file with upstream",
		ValidArgsFunction: completeFname,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := service.Sync(
				cmd.Context(),
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				getIndexFsProvider(),
				yes,
				args[0],
			); err != nil {
				if errors.Is(err, service.ErrCancelled) {
					os.Exit(1)
				}

				return fmt.Errorf("cannot sync '%s': %w", args[0], err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Automatic 'yes' to prompts")

	// cmd.Flags().BoolVarP(&patch, "patch", "p", false, "Interactively apply changes per entry")

	return cmd
}
