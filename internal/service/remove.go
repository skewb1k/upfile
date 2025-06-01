package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

func FileExistsAndReadable(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	_ = f.Close()
	return nil
}

func (s Service) Remove(
	ctx context.Context,
	path string,
) error {
	if err := FileExistsAndReadable(path); err != nil {
		return err
	}

	fname := filepath.Base(path)
	entryDir := filepath.Dir(path)

	if err := s.store.DeleteEntry(ctx, fname, entryDir); err != nil {
		return fmt.Errorf("delete entry: %w", err)
	}

	return nil
}
