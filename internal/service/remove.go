package service

import (
	"context"
	"fmt"
	"path/filepath"
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
		return fmt.Errorf("delete entry: %w", err)
	}

	return nil
}
