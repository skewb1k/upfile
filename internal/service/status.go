package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
)

func IsRealDirectory(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrDirNotExists
		}

		return err
	}

	if !info.IsDir() {
		return ErrNotDirectory
	}

	return nil
}

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

	rendered := make([]Entry, len(files))
	maxWidth := 0

	for i, fname := range files {
		rendered[i] = Entry{
			Path:   filepath.Join(dir, fname),
			Status: EntryStatusUpToDate,
			Err:    nil,
		}

		upstream, err := indexProvider.GetUpstream(ctx, fname)
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

		mustFmt(fmt.Fprintln(stdout, line))
	}

	return nil
}
