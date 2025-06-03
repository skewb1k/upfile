package cmd

import (
	"upfile/internal/service"

	"github.com/spf13/cobra"
)

func list() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List tracked files",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		RunE: wrap(func(cmd *cobra.Command, s *service.Service, args []string) error {
			files, err := s.List(cmd.Context())
			if err != nil {
				return err //nolint: wrapcheck
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
		}),
	}
}
