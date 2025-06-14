package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/charmbracelet/lipgloss"
	"github.com/skewb1k/upfile/internal/index"
	"github.com/skewb1k/upfile/pkg/validfname"
)

func List(
	ctx context.Context,
	stdout io.Writer,
	indexProvider IndexProvider,
	files []string,
) error {
	if len(files) == 0 {
		var err error
		files, err = indexProvider.GetFilenames(ctx)
		if err != nil {
			return fmt.Errorf("get files: %w", err)
		}
	} else {
		sort.Strings(files)
	}

	type upstreamFile struct {
		Upstream
		fname string
	}

	upstreamFiles := make([]upstreamFile, len(files))

	for fileIdx, fname := range files {
		if !validfname.ValidateFilename(fname) {
			return ErrInvalidFilename
		}

		upstream, err := indexProvider.GetUpstream(ctx, fname)
		if err != nil {
			if errors.Is(err, index.ErrNotFound) {
				return ErrNotTracked
			}

			return fmt.Errorf("get upstream: %w", err)
		}

		upstreamFiles[fileIdx] = upstreamFile{
			Upstream: upstream,
			fname:    fname,
		}
	}

	for fileIdx, upstream := range upstreamFiles {
		entriesList, err := indexProvider.GetEntriesByFilename(ctx, upstream.fname)
		if err != nil {
			return fmt.Errorf("get entries by filename: %w", err)
		}

		mustFmt(fmt.Println(_headingStyle.Render(upstream.fname)))

		renderedEntries := make([]Entry, len(entriesList))
		maxWidth := 0

		for entryIdx, entry := range entriesList {
			renderedEntries[entryIdx] = Entry{
				Path:   filepath.Join(entry, upstream.fname),
				Status: EntryStatusUpToDate,
				Err:    nil,
			}

			existing, err := os.ReadFile(renderedEntries[entryIdx].Path)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					renderedEntries[entryIdx].Status = EntryStatusDeleted
				} else {
					renderedEntries[entryIdx].Err = errors.Unwrap(err)
				}
			} else if !upstream.Hash.EqualBytes(existing) {
				renderedEntries[entryIdx].Status = EntryStatusModified
			}

			if w := len(renderedEntries[entryIdx].Path); w > maxWidth {
				maxWidth = w
			}
		}

		for _, e := range renderedEntries {
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

			mustFmt(fmt.Fprintln(stdout, line))
		}

		if fileIdx < len(files)-1 {
			mustFmt(fmt.Fprintln(stdout))
		}
	}

	return nil
}
