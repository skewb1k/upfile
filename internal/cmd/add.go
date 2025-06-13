package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/store"
	"github.com/spf13/cobra"
)

func addCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <path>",
		Short: "Add a file to be tracked",
		Args:  cobra.ExactArgs(1),
		RunE: withStore(func(cmd *cobra.Command, s *store.Store, args []string) error {
			path, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("failed to get abs path to file: %w", err)
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			entry, fname := filepath.Dir(path), filepath.Base(path)

			if err := s.CreateEntry(cmd.Context(), fname, entry); err != nil {
				if errors.Is(err, store.ErrExists) {
					cmd.Printf("File '%s' already tracked\n", path)
					return nil
				}

				return fmt.Errorf("create entry: %w", err)
			}

			upstreamExists, err := s.CheckUpstream(cmd.Context(), fname)
			if err != nil {
				return fmt.Errorf("check upstream: %w", err)
			}

			if !upstreamExists {
				if err := s.SetUpstream(cmd.Context(), fname, store.NewUpstream(content)); err != nil {
					return fmt.Errorf("set upstream: %w", err)
				}
			}

			return nil
		}),
	}

	return cmd
}
