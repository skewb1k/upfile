package cmd

import (
	"github.com/spf13/cobra"
)

const Version = "v0.0.1"

func version() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Println(Version)
			return nil
		},
	}
}
