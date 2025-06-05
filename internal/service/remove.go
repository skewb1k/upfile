package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"upfile/internal/index"
)

func (s Service) Remove(
	ctx context.Context,
	path string,
) error {
	exists, err := s.userfileProvider.CheckFile(ctx, path)
	if err != nil {
		return fmt.Errorf("check file: %w", err)
	}

	if !exists {
		return ErrFileNotFound
	}

	fname := filepath.Base(path)
	entryDir := filepath.Dir(path)

	if err := s.indexProvider.DeleteEntry(ctx, fname, entryDir); err != nil {
		if errors.Is(err, index.ErrNotFound) {
			return ErrNotTracked
		}

		return fmt.Errorf("delete entry: %w", err)
	}

	return nil
}
