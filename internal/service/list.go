package service

import (
	"context"
	"fmt"
)

type File struct {
	Fname   string
	Entries []string
}

func (s Service) List(
	ctx context.Context,
) ([]File, error) {
	files, err := s.indexProvider.GetFiles(ctx)
	if err != nil {
		return nil, fmt.Errorf("get files: %w", err)
	}

	res := make([]File, len(files))

	for i, fname := range files {
		entries, err := s.indexProvider.GetEntriesByFname(ctx, fname)
		if err != nil {
			return nil, fmt.Errorf("get entries by filename: %w", err)
		}

		res[i] = File{
			Fname:   fname,
			Entries: entries,
		}
	}

	return res, nil
}
