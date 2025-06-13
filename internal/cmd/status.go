package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/skewb1k/upfile/internal/store"
	"github.com/spf13/cobra"
)

func statusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status [<dir>]",
		Short: "Print status of files in dir (default: current dir)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := "."
			if len(args) == 1 {
				dir = args[0]
			}

			s := store.New(getBaseDir())

			absDir, err := filepath.Abs(dir)
			if err != nil {
				return fmt.Errorf("failed to get abs path to dir: %w", err)
			}

			files, err := s.GetFilenamesByEntry(cmd.Context(), absDir)
			if err != nil {
				if errors.Is(err, store.ErrNotFound) {
					if dir == "." {
						cmd.Println("No tracked files in this directory")
					} else {
						cmd.Printf("No tracked files in '%s'\n", dir)
					}

					return nil
				}

				return fmt.Errorf("get files by entry dir: %w", err)
			}

			rendered := make([]Entry, len(files))
			maxWidth := 0

			for i, fname := range files {
				rendered[i] = Entry{
					Path:   filepath.Join(absDir, fname),
					Status: EntryStatusUpToDate,
					Err:    nil,
				}

				upstream, err := s.GetUpstream(cmd.Context(), fname)
				if err != nil {
					return fmt.Errorf("get upstream: %w", err)
				}

				existing, err := os.ReadFile(rendered[i].Path)
				if err != nil {
					if errors.Is(err, os.ErrNotExist) {
						rendered[i].Status = EntryStatusDeleted
					} else {
						rendered[i].Err = errors.Unwrap(err)
					}
				} else if !upstream.Hash.EqualBytes(existing) {
					rendered[i].Status = EntryStatusModified
				}

				if w := len(rendered[i].Path); w > maxWidth {
					maxWidth = w
				}
			}

			for _, e := range rendered {
				var text string
				if e.Err != nil {
					text = "error: " + e.Err.Error()
				} else {
					text = statusAsString(e.Status)
				}

				line := lipgloss.JoinHorizontal(
					lipgloss.Top,
					lipgloss.NewStyle().Width(maxWidth+margin).Render(e.Path),
					text,
				)

				if e.Err != nil {
					line = _errorLineStyle.Render(line)
				}

				cmd.Println(line)
			}

			return nil
		},
	}

	return cmd
}
