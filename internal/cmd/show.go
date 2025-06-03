package cmd

import (
	"path/filepath"

	"upfile/internal/service"

	"github.com/spf13/cobra"
)

func show() *cobra.Command {
	return &cobra.Command{
		Use:   "show <filename>",
		Short: "Show upstream version of file",
		Args:  cobra.ExactArgs(1),
		RunE: wrap(func(cmd *cobra.Command, s *service.Service, args []string) error {
			fname := filepath.Base(args[0])

			content, err := s.Show(cmd.Context(), fname)
			if err != nil {
				return err //nolint: wrapcheck
			}

			cmd.Print(content)
			return nil
		}),
	}
}
