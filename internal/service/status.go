package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type EntryStatus int

const (
	EntryStatusModified EntryStatus = iota
	EntryStatusUpToDate
	EntryStatusDeleted
)

type Entry struct {
	Fname  string
	Status EntryStatus
}

func (s Service) Status(
	ctx context.Context,
	dir string,
) ([]Entry, error) {
	files, err := s.indexProvider.GetFilesByEntryDir(ctx, dir)
	if err != nil {
		return nil, fmt.Errorf("get files by entry dir: %w", err)
	}

	res := make([]Entry, len(files))

	for i, filename := range files {
		res[i] = Entry{
			Fname:  filename,
			Status: EntryStatusUpToDate,
		}

		upstream, err := s.indexProvider.GetUpstream(ctx, filename)
		if err != nil {
			return nil, fmt.Errorf("get upstream: %w", err)
		}

		existing, err := os.ReadFile(filepath.Join(dir, filename))
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return nil, fmt.Errorf("read file: %w", err)
			}

			res[i].Status = EntryStatusDeleted
		} else {
			existingHash := sha256.Sum256(existing)
			upstreamHash := sha256.Sum256([]byte(upstream))

			if existingHash != upstreamHash {
				res[i].Status = EntryStatusModified
			}
		}
	}

	return res, nil
}
