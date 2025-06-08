package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/userfile"
)

type File struct {
	Fname   string
	Entries []Entry
}

func (s Service) List(
	ctx context.Context,
) ([]File, error) {
	files, err := s.indexProvider.GetFilenames(ctx)
	if err != nil {
		return nil, fmt.Errorf("get files: %w", err)
	}

	res := make([]File, len(files))

	// TODO: refactor, do not duplicate Status command logic
	for i, fname := range files {
		entries, err := s.indexProvider.GetEntriesByFilename(ctx, fname)
		if err != nil {
			return nil, fmt.Errorf("get entries by filename: %w", err)
		}

		upstream, err := s.indexProvider.GetUpstream(ctx, fname)
		if err != nil {
			return nil, fmt.Errorf("get upstream: %w", err)
		}

		e := make([]Entry, len(entries))

		for j, entry := range entries {
			path := filepath.Join(entry, fname)
			e[j] = Entry{
				Path:   path,
				Status: EntryStatusUpToDate,
			}

			existing, err := s.userfileProvider.ReadFile(ctx, path)
			if err != nil {
				if !errors.Is(err, userfile.ErrNotFound) {
					return nil, fmt.Errorf("read file: %w", err)
				}

				e[j].Status = EntryStatusDeleted
			} else if !upstream.Hash.EqualString(existing) {
				e[j].Status = EntryStatusModified
			}
		}

		res[i] = File{
			Fname:   fname,
			Entries: e,
		}
	}

	return res, nil
}

func (s Service) ListTrackedFilenames(ctx context.Context) ([]string, error) {
	files, err := s.indexProvider.GetFilenames(ctx)
	if err != nil {
		return nil, fmt.Errorf("get files: %w", err)
	}

	return files, nil
}
