package cmd

import (
	"github.com/spf13/cobra"
)

func sync() *cobra.Command {
	// TODO:
	// var patch bool

	cmd := &cobra.Command{
		Use:   "sync <filename>",
		Short: "Sync all outdated entries of file with upstream",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// s := service.New(storeFs.New(getBaseDir()))
			// err := s.Sync(cmd.Context(), args[0])
			// if err != nil {
			// 	return err
			// }

			return nil
		},
	}

	// cmd.Flags().BoolVarP(&patch, "patch", "p", false, "Interactively apply changes per entry")

	return cmd
}
