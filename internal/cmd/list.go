package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/skewb1k/upfile/internal/entries"
	"github.com/skewb1k/upfile/internal/upstreams"
	"github.com/spf13/cobra"
)

type EntryStatus int

const (
	EntryStatusModified EntryStatus = iota
	EntryStatusUpToDate
	EntryStatusDeleted
)

type Entry struct {
	Path   string
	Status EntryStatus
	Err    error
}

type File struct {
	Fname   string
	Entries []Entry
}

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List tracked files",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			baseDir := getBaseDir()
			upstreamsProvider := upstreams.NewProvider(baseDir)
			entriesProvider := entries.NewProvider(baseDir)

			files, err := upstreamsProvider.GetFilenames(cmd.Context())
			if err != nil {
				return fmt.Errorf("get files: %w", err)
			}

			res := make([]File, len(files))

			// TODO: refactor, do not duplicate Status command logic
			for i, fname := range files {
				upstream, err := upstreamsProvider.GetUpstream(cmd.Context(), fname)
				if err != nil {
					return fmt.Errorf("get upstream: %w", err)
				}

				entriesList, err := entriesProvider.GetEntriesByFilename(cmd.Context(), fname)
				if err != nil {
					return fmt.Errorf("get entries by filename: %w", err)
				}

				res[i] = File{
					Fname:   fname,
					Entries: make([]Entry, len(entriesList)),
				}

				for j, entry := range entriesList {
					path := filepath.Join(entry, fname)
					res[i].Entries[j] = Entry{
						Path:   path,
						Status: EntryStatusUpToDate,
						Err:    nil,
					}

					existing, err := os.ReadFile(path)
					if err != nil {
						if errors.Is(err, os.ErrNotExist) {
							res[i].Entries[j].Status = EntryStatusDeleted
						} else {
							res[i].Entries[j].Err = err
						}
					} else if !upstream.Hash.EqualBytes(existing) {
						res[i].Entries[j].Status = EntryStatusModified
					}
				}
			}

			if len(files) == 0 {
				return nil
			}

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			defer w.Flush()

			for i, f := range res {
				mustFprintf(w, "%s:\n", f.Fname)

				// TODO: fix alignment and refactor
				for _, entry := range f.Entries {
					fn := entry.Path
					if entry.Err != nil {
						mustFprintf(w, red("\t%s\t%s\n"), fn, errors.Unwrap(entry.Err).Error())
					} else {
						mustFprintf(w, "\t%s\t%s\n", fn, statusAsString(entry.Status))
					}
				}

				if i < len(files)-1 {
					mustFprintf(w, "\n")
				}
			}

			return nil
		},
	}

	return cmd
}
