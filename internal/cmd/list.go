package cmd

import (
	"fmt"

	"github.com/skewb1k/upfile/internal/service"
	"github.com/spf13/cobra"
)

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "list [<filename>...]",
		Short:             "List tracked files and entries status",
		Aliases:           []string{"ls"},
		Args:              cobra.ArbitraryArgs,
		ValidArgsFunction: completeFnames,
		RunE: func(cmd *cobra.Command, files []string) error {
			if err := service.List(cmd.Context(), cmd.OutOrStdout(), getIndexFsProvider(), files); err != nil {
				return fmt.Errorf("cannot list: %w", err)
			}

			return nil
		},
	}

	return cmd
}
