package commands

import (
	"fmt"

	"github.com/skewb1k/upfile/internal/service"
	"github.com/spf13/cobra"
)

func List() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [<filename>...]",
		Short: "List tracked files and entries status",
		Long: doc(`
List tracked files and entries status.

If no filenames are provided, all tracked files will be listed.

If one or more filenames are specified, only those files will be listed.
`),
		Aliases:           []string{"ls"},
		Args:              cobra.ArbitraryArgs,
		ValidArgsFunction: completeFnames,
		RunE: func(cmd *cobra.Command, files []string) error {
			if err := service.List(cmd.Context(), cmd.OutOrStdout(), getIndexFsProvider(), files); err != nil {
				// FIXME:
				return fmt.Errorf("%w", err)
			}

			return nil
		},
	}

	return cmd
}
