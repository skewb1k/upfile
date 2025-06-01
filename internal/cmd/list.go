package cmd

import (
	"upfile/internal/service"
	storeFs "upfile/internal/store/fs"

	"github.com/spf13/cobra"
)

func list() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List tracked files",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			s := service.New(storeFs.New(getBaseDir()))

			files, err := s.List(cmd.Context())
			if err != nil {
				return err
			}

			if len(files) == 0 {
				return nil
			}

			for _, f := range files {
				cmd.Println(f.Fname)
				for _, dir := range f.Entries {
					cmd.Printf("  %s\n", dir)
				}
				cmd.Println()
			}

			return nil
		},
	}
}
