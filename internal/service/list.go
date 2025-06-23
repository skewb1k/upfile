package service

import (
	"context"
	"errors"
	"fmt"
	"io"
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

	for i, upstream := range upstreamFiles {
		entriesList, err := indexProvider.GetEntriesByFilename(ctx, upstream.fname)
		if err != nil {
			return fmt.Errorf("get entries by filename: %w", err)
		}

		mustFmt(fmt.Println(_headingStyle.Render(upstream.fname)))

		renderedEntries := make([]*Entry, len(entriesList))
		maxWidth := 0

		for entryIdx, entry := range entriesList {
			path := filepath.Join(entry, upstream.fname)

			renderedEntries[entryIdx] = getEntry(path, upstream.Hash)

			if w := len(path); w > maxWidth {
				maxWidth = w
			}
		}

		for _, e := range renderedEntries {
			var text string
			if e.Err != nil {
				text = "error: " + e.Err.Error()
			} else {
				text = e.Status
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

		if i < len(files)-1 {
			mustFmt(fmt.Fprintln(stdout))
		}
	}

	return nil
}
