package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/userfile"
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
}

func (s Service) Status(
	ctx context.Context,
	dir string,
) ([]Entry, error) {
	files, err := s.indexProvider.GetFilenamesByEntry(ctx, dir)
	if err != nil {
		return nil, fmt.Errorf("get files by entry dir: %w", err)
	}

	res := make([]Entry, len(files))

	for i, fname := range files {
		path := filepath.Join(dir, fname)
		res[i] = Entry{
			Path:   path,
			Status: EntryStatusUpToDate,
		}

		upstream, err := s.indexProvider.GetUpstream(ctx, fname)
		if err != nil {
			return nil, fmt.Errorf("get upstream: %w", err)
		}

		existing, err := s.userfileProvider.ReadFile(ctx, path)
		if err != nil {
			if !errors.Is(err, userfile.ErrNotFound) {
				return nil, fmt.Errorf("read file: %w", err)
			}

			res[i].Status = EntryStatusDeleted
		} else if !upstream.Hash.EqualString(existing) {
			res[i].Status = EntryStatusModified
		}
	}

	return res, nil
}
