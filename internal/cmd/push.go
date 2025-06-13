package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/store"
	"github.com/spf13/cobra"
)

func pushCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "push <path>",
		Short: "Push file to the upstream",
		Args:  cobra.ExactArgs(1),
		RunE: withStore(func(cmd *cobra.Command, s *store.Store, args []string) error {
			path, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("failed to get abs path to file: %w", err)
			}

			newContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			entry, fname := filepath.Dir(path), filepath.Base(path)

			exists, err := s.CheckEntry(cmd.Context(), fname, entry)
			if err != nil {
				return fmt.Errorf("check entry: %w", err)
			}

			if !exists {
				return ErrNotTracked
			}

			upstream, err := s.GetUpstream(cmd.Context(), fname)
			if err != nil {
				return fmt.Errorf("get upstream: %w", err)
			}

			if upstream.Hash.EqualBytes(newContent) {
				cmd.Println("File up-to-date")
				return nil
			}

			if err := s.SetUpstream(
				cmd.Context(),
				fname,
				store.NewUpstream(newContent),
			); err != nil {
				return fmt.Errorf("set upstream: %w", err)
			}

			return nil
		}),
	}
}
