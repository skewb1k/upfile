package cmd

import (
	"strings"

	"upfile/internal/service"
	storeFs "upfile/internal/store/fs"

	"github.com/spf13/cobra"
)

func list() *cobra.Command {
	return &cobra.Command{
		Use:   "list <filename>",
		Short: "List tracked files",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			baseDir := getBaseDir()

			files, err := service.GetFiles(cmd.Context(),
				storeFs.New(baseDir),
			)
			if err != nil {
				return err
			}

			cmd.Println(strings.Join(files, "\n"))

			return nil
		},
	}
}
