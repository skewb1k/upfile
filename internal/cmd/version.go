package cmd

import (
	"github.com/spf13/cobra"
)

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "version",
		Aliases:           []string{"v"},
		Short:             "Print version",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Println(version)

			return nil
		},
	}
}
