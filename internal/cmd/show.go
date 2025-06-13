package cmd

import (
	"errors"
	"fmt"

	"github.com/skewb1k/upfile/internal/store"
	"github.com/spf13/cobra"
)

func showCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "show <filename>",
		Short:             "Show upstream version of file",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completeFname,
		RunE: withStore(func(cmd *cobra.Command, s *store.Store, args []string) error {
			upstream, err := s.GetUpstream(cmd.Context(), args[0])
			if err != nil {
				if errors.Is(err, store.ErrNotFound) {
					return ErrNotTracked
				}

				return fmt.Errorf("get upstream: %w", err)
			}

			cmd.Print(upstream.Content)

			return nil
		}),
	}

	return cmd
}
