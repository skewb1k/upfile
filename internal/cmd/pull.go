package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/store"
	"github.com/spf13/cobra"
)

func pullCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "pull <path>",
		Short: "Pull file from upstream",
		Args:  cobra.ExactArgs(1),
		RunE: withStore(func(cmd *cobra.Command, s *store.Store, args []string) error {
			path, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("failed to get abs path to dest dir: %w", err)
			}

			destDir, fname := filepath.Dir(path), filepath.Base(path)

			upstream, err := s.GetUpstream(cmd.Context(), fname)
			if err != nil {
				if errors.Is(err, store.ErrNotFound) {
					return ErrNotTracked
				}

				return fmt.Errorf("get upstream: %w", err)
			}

			existing, err := os.ReadFile(path)
			if err != nil {
				if !errors.Is(err, os.ErrNotExist) {
					return err
				}
			} else if upstream.Hash.EqualBytes(existing) {
				cmd.Println("File up-to-date")
				return nil
			}

			if !yes {
				cmd.Printf("File '%s' will be updated\n", path)

				ok, err := askDefaultYes(cmd.InOrStdin(), cmd.OutOrStdout())
				if err != nil {
					return err
				}

				if !ok {
					os.Exit(1)
				}
			}

			if err := s.CreateEntry(cmd.Context(), fname, destDir); err != nil {
				if !errors.Is(err, store.ErrExists) {
					return fmt.Errorf("create entry: %w", err)
				}
			}

			if err := MkdirAllWriteFile(path, upstream.Content); err != nil {
				return err
			}

			return nil
		}),
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Automatic 'yes' to prompts")

	return cmd
}
