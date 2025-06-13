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

const margin = 2

var (
	_headingStyle = lipgloss.NewStyle().Bold(true)
	_pathStyle    = lipgloss.NewStyle().MarginLeft(margin)
)

var _errorLineStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))

func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List tracked files",
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			s := store.New(getBaseDir())

			files, err := s.GetFilenames(cmd.Context())
			if err != nil {
				return fmt.Errorf("get files: %w", err)
			}

			for i, fname := range files {
				upstream, err := s.GetUpstream(cmd.Context(), fname)
				if err != nil {
					panic(fmt.Errorf("get upstream: %w", err))
				}

				entriesList, err := s.GetEntriesByFilename(cmd.Context(), fname)
				if err != nil {
					return fmt.Errorf("get entries by filename: %w", err)
				}

				fmt.Println(_headingStyle.Render(fname))

				rendered := make([]Entry, len(entriesList))
				maxWidth := 0

				for j, entry := range entriesList {
					rendered[j] = Entry{
						Path:   filepath.Join(entry, fname),
						Status: EntryStatusUpToDate,
						Err:    nil,
					}

					existing, err := os.ReadFile(rendered[j].Path)
					if err != nil {
						if errors.Is(err, os.ErrNotExist) {
							rendered[j].Status = EntryStatusDeleted
						} else {
							rendered[j].Err = errors.Unwrap(err)
						}
					} else if !upstream.Hash.EqualBytes(existing) {
						rendered[j].Status = EntryStatusModified
					}

					if w := len(rendered[j].Path); w > maxWidth {
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
						_pathStyle.Width(maxWidth+margin).Render(e.Path),
						text,
					)

					if e.Err != nil {
						line = _errorLineStyle.Render(line)
					}

					cmd.Println(line)
				}

				if i < len(files)-1 {
					cmd.Println()
				}
			}

			return nil
		},
	}

	return cmd
}

var (
	_modifiedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	_upToDateStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	_deletedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

func statusAsString(status EntryStatus) string {
	switch status {
	case EntryStatusModified:
		return _modifiedStyle.Render("Modified")
	case EntryStatusUpToDate:
		return _upToDateStyle.Render("Up-to-date")
	case EntryStatusDeleted:
		return _deletedStyle.Render("Deleted")
	default:
		panic("UNEXPECTED")
	}
}
