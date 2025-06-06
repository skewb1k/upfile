package cmd

import (
	"fmt"
	"path/filepath"
	"text/tabwriter"

	"github.com/skewb1k/upfile/internal/service"

	"github.com/spf13/cobra"
)

func status() *cobra.Command {
	return &cobra.Command{
		Use:   "status [<dir>]",
		Short: "Print status of files in dir (default: current dir)",
		Args:  cobra.MaximumNArgs(1),
		RunE: wrap(func(cmd *cobra.Command, s *service.Service, args []string) error {
			dir := "."
			if len(args) == 1 {
				dir = args[0]
			}

			absDir, err := filepath.Abs(dir)
			if err != nil {
				return fmt.Errorf("failed to get abs path to dir: %w", err)
			}

			entries, err := s.Status(cmd.Context(), absDir)
			if err != nil {
				return err //nolint: wrapcheck
			}

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			defer w.Flush()

			for _, entry := range entries {
				mustFprintf(w, "%s\t%s\n", filepath.Base(entry.Path), statusAsString(entry.Status))
			}

			return nil
		}),
	}
}
