package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "v0.0.1"

func version() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\n", Version)
			if err != nil {
				return fmt.Errorf("failed to print: %w", err)
			}

			return nil
		},
	}
}
