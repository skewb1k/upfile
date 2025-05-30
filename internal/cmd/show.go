package cmd

import (
	"github.com/spf13/cobra"
)

func show() *cobra.Command {
	return &cobra.Command{
		Use:   "show <filename>",
		Short: "Show upstream version of file",
		Args:  cobra.ExactArgs(1),
		// RunE: func(cmd *cobra.Command, args []string) error {
		// 	baseDir := getBaseDir()
		//
		// 	upstreamContent, err := service.GetUpstream(cmd.Context(),
		// 		commitsFs.New(baseDir),
		// 		headFs.New(baseDir),
		// 		args[0],
		// 	)
		// 	if err != nil {
		// 		return err
		// 	}
		//
		// 	cmd.Print(upstreamContent)
		//
		// 	return nil
		// },
	}
}
