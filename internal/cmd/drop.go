package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/skewb1k/upfile/internal/service"
	"github.com/spf13/cobra"
)

func dropCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "drop <filename>",
		Args:  cobra.ExactArgs(1),
		Short: "Remove tracked file upstream and entries",
		Long: doc(`
Permanently removes a file from UpFile tracking:

- Deletes the upstream version
- Removes all tracked entries from the index

Note: This does NOT delete any actual files from user-space filesystem.

Use with caution. You will be prompted to confirm removal unless --yes is specified.
`),
		ValidArgsFunction: completeFname,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := service.Drop(
				cmd.Context(),
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				getStore(),
				yes,
				args[0],
			); err != nil {
				if errors.Is(err, service.ErrCancelled) {
					os.Exit(1)
				}

				return fmt.Errorf("cannot drop '%s': %w", args[0], err)
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Automatic 'yes' to prompts")

	return cmd
}
