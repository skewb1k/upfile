package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/service"
	"github.com/spf13/cobra"
)

func pushCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "push <path>...",
		Short:             "Push tracked file(s) to the upstream",
		Args:              cobra.MinimumNArgs(1),
		ValidArgsFunction: completeEntry,
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				path, err := filepath.Abs(arg)
				if err != nil {
					return fmt.Errorf("failed to get absolute path to file '%s': %w", arg, err)
				}

				if err := service.Push(
					cmd.Context(),
					cmd.OutOrStdout(),
					getIndexFsProvider(),
					path,
				); err != nil {
					return fmt.Errorf("cannot push '%s': %w", path, err)
				}
			}
			return nil
		},
	}

	return cmd
}
