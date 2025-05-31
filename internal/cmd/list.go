package cmd

import (
	"strings"

	"upfile/internal/service"
	storeFs "upfile/internal/store/fs"

	"github.com/spf13/cobra"
)

func list() *cobra.Command {
	return &cobra.Command{
		Use:     "list <filename>",
		Short:   "List tracked files",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			s := service.New(storeFs.New(getBaseDir()))

			files, err := s.GetFiles(cmd.Context())
			if err != nil {
				return err
			}

			cmd.Println(strings.Join(files, "\n"))

			return nil
		},
	}
}
