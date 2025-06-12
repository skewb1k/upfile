package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/skewb1k/upfile/internal/service"
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
		RunE: wrap(func(cmd *cobra.Command, s *service.Service, args []string) error {
			files, err := s.GetFilenames(cmd.Context())
			if err != nil {
				return fmt.Errorf("get files: %w", err)
			}

			res := make([]File, len(files))

			// TODO: refactor, do not duplicate Status command logic
			for i, fname := range files {
				head, err := s.GetHead(cmd.Context(), fname)
				if err != nil {
					return fmt.Errorf("get head: %w", err)
				}

				headCommit, err := s.GetCommit(cmd.Context(), head)
				if err != nil {
					return fmt.Errorf("get commit: %w", err)
				}

				entriesList, err := s.GetEntriesByFilename(cmd.Context(), fname)
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
					} else if !headCommit.ContentHash.EqualBytes(existing) {
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
					if entry.Err != nil {
						mustFprintf(w, red("\t%s\t%s\n"), entry.Path, errors.Unwrap(entry.Err).Error())
					} else {
						mustFprintf(w, "\t%s\t%s\n", entry.Path, statusAsString(entry.Status))
					}
				}

				if i < len(files)-1 {
					mustFprintf(w, "\n")
				}
			}

			return nil
		}),
	}

	return cmd
}
