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
		Use:   "drop <filename>...",
		Short: "Remove tracked file(s) upstream and entries",
		Long: doc(`
Permanently removes one or more files from UpFile tracking:

- Deletes the upstream version
- Removes all tracked entries from the index

Note: This does NOT delete any actual files from the user-space filesystem.

Use with caution. You will be prompted to confirm removal unless --yes is specified.
`),
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: completeFnames,
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, name := range args {
				if err := service.Drop(
					cmd.Context(),
					cmd.InOrStdin(),
					cmd.OutOrStdout(),
					getIndexFsProvider(),
					yes,
					name,
				); err != nil {
					if errors.Is(err, service.ErrCancelled) {
						os.Exit(1)
					}
					return fmt.Errorf("cannot drop '%s': %w", name, err)
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Automatic 'yes' to prompts")

	return cmd
}
