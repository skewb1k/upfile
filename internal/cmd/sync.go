package cmd

import (
	"path/filepath"

	"upfile/internal/service"

	"github.com/spf13/cobra"
)

func sync() *cobra.Command {
	// TODO:
	// var patch bool

	cmd := &cobra.Command{
		Use:   "sync <filename>",
		Short: "Sync all entries of file with upstream",
		Args:  cobra.ExactArgs(1),
		RunE: wrap(func(cmd *cobra.Command, s *service.Service, args []string) error {
			fname := filepath.Base(args[0])

			if err := s.Sync(cmd.Context(), fname); err != nil {
				return err //nolint: wrapcheck
			}

			return nil
		}),
	}

	// cmd.Flags().BoolVarP(&patch, "patch", "p", false, "Interactively apply changes per entry")

	return cmd
}
