package cmd

import (
	"text/tabwriter"

	"github.com/skewb1k/upfile/internal/service"

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

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			defer w.Flush()

			for i, f := range files {
				mustFprintf(w, "%s:\n", f.Fname)

				for _, entry := range f.Entries {
					mustFprintf(w, "\t%s\t%s\n", entry.Path, statusAsString(entry.Status))
				}

				if i < len(files)-1 {
					mustFprintf(w, "\n")
				}
			}

			return nil
		}),
	}
}
