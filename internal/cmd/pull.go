package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"upfile/internal/service"

	"github.com/spf13/cobra"
)

func pull() *cobra.Command {
	var dest string

	cmd := &cobra.Command{
		Use:   "pull <filename>",
		Short: "Pull file from origin",
		Args:  cobra.ExactArgs(1),
		RunE: wrap(func(cmd *cobra.Command, s *service.Service, args []string) error {
			destAbs, err := filepath.Abs(dest)
			if err != nil {
				return fmt.Errorf("failed to get abs path to dest dir: %w", err)
			}

			fname := filepath.Base(args[0])

			if err := s.Pull(cmd.Context(), fname, destAbs); err != nil {
				if errors.Is(err, service.ErrUpToDate) {
					cmd.Println("File up-to-date")
					return nil
				}

				return err //nolint: wrapcheck
			}

			return nil
		}),
	}

	cmd.Flags().StringVarP(&dest, "dest", "d", ".", "Destination folder")

	return cmd
}
