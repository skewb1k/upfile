package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/store"
	"github.com/spf13/cobra"
)

func syncCmd() *cobra.Command {
	// TODO:
	// like git add -p
	// var patch bool

	var yes bool

	cmd := &cobra.Command{
		Use:               "sync <filename>",
		Args:              cobra.ExactArgs(1),
		Short:             "Sync all entries of file with upstream",
		ValidArgsFunction: completeFname,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := store.New(getBaseDir())

			fname := args[0]

			entries, err := s.GetEntriesByFilename(cmd.Context(), fname)
			if err != nil {
				return fmt.Errorf("get entries by filename: %w", err)
			}

			upstream, err := s.GetUpstream(cmd.Context(), fname)
			if err != nil {
				if errors.Is(err, store.ErrNotFound) {
					return ErrNotTracked
				}

				return fmt.Errorf("get upstream: %w", err)
			}

			toUpdate := make([]string, 0)

			for _, entry := range entries {
				path := filepath.Join(entry, fname)

				existing, err := os.ReadFile(path)
				if err == nil && upstream.Hash.EqualBytes(existing) {
					// Up-to-date, skip
					continue
				}

				if err != nil && !errors.Is(err, os.ErrNotExist) {
					return err
				}

				toUpdate = append(toUpdate, path)
			}

			if len(toUpdate) == 0 {
				cmd.Println("Everything up-to-date")
				return nil
			}

			if !yes && !ask(cmd.InOrStdin(), toUpdate, true /* defaultYes */, "The following entries will be updated:") {
				os.Exit(1)
				return nil
			}

			for _, fullPath := range toUpdate {
				if err := WriteFile(fullPath, upstream.Content); err != nil {
					return fmt.Errorf("write file: %w", err)
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Automatic 'yes' to prompts")

	// cmd.Flags().BoolVarP(&patch, "patch", "p", false, "Interactively apply changes per entry")

	return cmd
}
