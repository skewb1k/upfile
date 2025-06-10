package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/entries"
	"github.com/skewb1k/upfile/internal/upstreams"
	"github.com/spf13/cobra"
)

func syncCmd() *cobra.Command {
	// TODO:
	// var patch bool

	var yes bool

	cmd := &cobra.Command{
		Use:   "sync <filename>",
		Args:  cobra.ExactArgs(1),
		Short: "Sync all entries of file with upstream",

		RunE: func(cmd *cobra.Command, args []string) error {
			baseDir := getBaseDir()
			upstreamsProvider := upstreams.NewProvider(baseDir)
			entriesProvider := entries.NewProvider(baseDir)

			fname := args[0]

			entriesList, err := entriesProvider.GetEntriesByFilename(cmd.Context(), fname)
			if err != nil {
				return fmt.Errorf("get entries by filename: %w", err)
			}

			upstream, err := upstreamsProvider.GetUpstream(cmd.Context(), fname)
			if err != nil {
				if errors.Is(err, upstreams.ErrNotFound) {
					return ErrNotTracked
				}

				return fmt.Errorf("get upstream: %w", err)
			}

			toUpdate := make([]string, 0)

			for _, entryDir := range entriesList {
				path := filepath.Join(entryDir, fname)

				existing, err := os.ReadFile(path)
				if err != nil {
					return err
				}

				if !upstream.Hash.EqualBytes(existing) {
					toUpdate = append(toUpdate, filepath.Join(entryDir, fname))
				}
			}

			if len(toUpdate) == 0 {
				return ErrUpToDate
			}

			if !yes && !ask(cmd.InOrStdin(), toUpdate, true, "The following entries will be updated:") {
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
	cmd.ValidArgsFunction = completeFname

	// cmd.Flags().BoolVarP(&patch, "patch", "p", false, "Interactively apply changes per entry")

	return cmd
}
