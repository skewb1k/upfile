package service

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
)

func Status(
	ctx context.Context,
	stdout io.Writer,
	indexProvider IndexProvider,
	dir string,
) error {
	if err := IsRealDirectory(dir); err != nil {
		return err
	}

	files, err := indexProvider.GetFilenamesByEntry(ctx, dir)
	if err != nil {
		return fmt.Errorf("get files by entry dir: %w", err)
	}

	if len(files) == 0 {
		mustFmt(fmt.Fprintf(stdout, "No tracked files in '%s'\n", dir))
		return nil
	}

	rendered := make([]*Entry, len(files))
	maxWidth := 0

	for i, fname := range files {
		upstream, err := indexProvider.GetUpstream(ctx, fname)
		if err != nil {
			return fmt.Errorf("get upstream: %w", err)
		}

		path := filepath.Join(dir, fname)

		rendered[i] = getEntry(path, upstream.Hash)

		if w := len(path); w > maxWidth {
			maxWidth = w
		}
	}

	for _, e := range rendered {
		var text string
		if e.Err != nil {
			text = "error: " + e.Err.Error()
		} else {
			text = e.Status
		}

		line := lipgloss.JoinHorizontal(
			lipgloss.Top,
			lipgloss.NewStyle().Width(maxWidth+margin).Render(e.Path),
			text,
		)

		if e.Err != nil {
			line = _errorLineStyle.Render(line)
		}

		mustFmt(fmt.Fprintln(stdout, line))
	}

	return nil
}
