package service

import (
	"context"
	"fmt"
	"path/filepath"
)

func (s Service) Sync(
	ctx context.Context,
	fname string,
) error {
	entries, err := s.indexProvider.GetEntriesByFname(ctx, fname)
	if err != nil {
		return fmt.Errorf("get entries by filename: %w", err)
	}

	upstream, err := s.indexProvider.GetUpstream(ctx, fname)
	if err != nil {
		return fmt.Errorf("get upstream: %w", err)
	}

	for _, entryDir := range entries {
		fullPath := filepath.Join(entryDir, fname)

		if err := s.userfileProvider.WriteFile(ctx, fullPath, upstream); err != nil {
			return fmt.Errorf("write file: %w", err)
		}
	}

	return nil
}
