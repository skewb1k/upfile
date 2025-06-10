package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/entries"
	"github.com/spf13/cobra"
)

func removeCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "remove <path>",
		Short:   "Remove entry from tracked list",
		Aliases: []string{"rm"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			baseDir := getBaseDir()
			entriesProvider := entries.NewProvider(baseDir)

			path, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("failed to get abs path to file: %w", err)
			}

			entryDir, fname := filepath.Dir(path), filepath.Base(path)

			if err := entriesProvider.DeleteEntry(cmd.Context(), fname, entryDir); err != nil {
				if errors.Is(err, entries.ErrNotFound) {
					return ErrNotTracked
				}

				return fmt.Errorf("delete entry: %w", err)
			}

			return nil
		},
	}
}
