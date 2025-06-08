package service

import (
	"context"
	"fmt"
	"path/filepath"
)

func (s Service) Sync(
	ctx context.Context,
	fname string,
	confirm func([]string) bool,
) error {
	entries, err := s.indexProvider.GetEntriesByFilename(ctx, fname)
	if err != nil {
		return fmt.Errorf("get entries by filename: %w", err)
	}

	upstream, err := s.indexProvider.GetUpstream(ctx, fname)
	if err != nil {
		return fmt.Errorf("get upstream: %w", err)
	}

	toUpdate := make([]string, 0)

	for _, entryDir := range entries {
		path := filepath.Join(entryDir, fname)

		existing, err := s.userfileProvider.ReadFile(ctx, path)
		if err != nil {
			return fmt.Errorf("read file: %w", err)
		}

		if !upstream.Hash.EqualString(existing) {
			toUpdate = append(toUpdate, filepath.Join(entryDir, fname))
		}
	}

	if len(toUpdate) == 0 {
		return ErrUpToDate
	}

	if !confirm(toUpdate) {
		return ErrCancelled
	}

	for _, fullPath := range toUpdate {
		if err := s.userfileProvider.WriteFile(ctx, fullPath, upstream.Content); err != nil {
			return fmt.Errorf("write file: %w", err)
		}
	}

	return nil
}
