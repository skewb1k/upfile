package cmd

import (
	"path/filepath"

	"upfile/internal/service"
	storeFs "upfile/internal/store/fs"

	"github.com/spf13/cobra"
)

func show() *cobra.Command {
	return &cobra.Command{
		Use:   "show <filename>",
		Short: "Show upstream version of file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s := service.New(storeFs.New(getBaseDir()))

			upstreamContent, err := s.Show(cmd.Context(), filepath.Base(args[0]))
			if err != nil {
				return err
			}

			cmd.Print(upstreamContent)

			return nil
		},
	}
}
