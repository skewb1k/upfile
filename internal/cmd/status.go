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

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status [<dir>]",
		Short: "Print status of files in dir (default: current dir)",
		Args:  cobra.MaximumNArgs(1),
		RunE: wrap(func(cmd *cobra.Command, s *service.Service, args []string) error {
			dir := "."
			if len(args) == 1 {
				dir = args[0]
			}

			absDir, err := filepath.Abs(dir)
			if err != nil {
				return fmt.Errorf("failed to get abs path to dir: %w", err)
			}

			files, err := s.GetFilenamesByEntry(cmd.Context(), absDir)
			if err != nil {
				// if errors.Is(err, index.ErrNotFound) {
				// 	return ErrNoEntries
				// }

				return fmt.Errorf("get files by entry dir: %w", err)
			}

			res := make([]Entry, len(files))

			for i, fname := range files {
				path := filepath.Join(absDir, fname)
				res[i] = Entry{
					Path:   path,
					Status: EntryStatusUpToDate,
					Err:    nil,
				}

				upstream, err := s.GetUpstream(cmd.Context(), fname)
				if err != nil {
					return fmt.Errorf("get upstream: %w", err)
				}

				existing, err := os.ReadFile(path)
				if err != nil {
					if !errors.Is(err, os.ErrNotExist) {
						return err
					}

					res[i].Status = EntryStatusDeleted
				} else if !upstream.Hash.EqualBytes(existing) {
					res[i].Status = EntryStatusModified
				}
			}

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			defer w.Flush()

			for _, entry := range res {
				mustFprintf(w, "%s\t%s\n", filepath.Base(entry.Path), statusAsString(entry.Status))
			}

			return nil
		}),
	}
}
