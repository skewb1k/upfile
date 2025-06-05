package service

import (
	"context"
	"errors"
	"fmt"

	"upfile/internal/index"
)

func (s Service) Drop(
	ctx context.Context,
	fname string,
	confirm func([]string) bool,
) error {
	// TODO: collect errors
	entries, err := s.indexProvider.GetEntriesByFname(ctx, fname)
	if err != nil {
		return fmt.Errorf("get entries by filename: %w", err)
	}

	if len(entries) == 0 {
		return ErrNoEntries
	}

	if !confirm(entries) {
		return ErrCancelled
	}

	if err := s.indexProvider.DeleteUpstream(ctx, fname); err != nil {
		if errors.Is(err, index.ErrNotFound) {
			return ErrNotTracked
		}

		return fmt.Errorf("delete upstream: %w", err)
	}

	for _, entry := range entries {
		if err := s.indexProvider.DeleteEntry(ctx, fname, entry); err != nil {
			return fmt.Errorf("delete entry: %w", err)
		}
	}

	return nil
}
