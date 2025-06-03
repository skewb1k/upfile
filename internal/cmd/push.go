package cmd

import (
	"errors"

	"upfile/internal/service"

	"github.com/spf13/cobra"
)

func push() *cobra.Command {
	return &cobra.Command{
		Use:   "push <path>",
		Short: "Push file to the upstream",
		Args:  cobra.ExactArgs(1),
		RunE: wrap(func(cmd *cobra.Command, s *service.Service, args []string) error {
			if err := s.Push(cmd.Context(), args[0]); err != nil {
				if errors.Is(err, service.ErrUpToDate) {
					cmd.Println("Already up to date")
					return nil
				}

				return err //nolint: wrapcheck
			}

			return nil
		}),
	}
}
