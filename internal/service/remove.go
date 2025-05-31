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
	entryPath string,
) error {
	if err := FileExistsAndReadable(entryPath); err != nil {
		return err
	}

	fname := filepath.Base(entryPath)
	entryDir := filepath.Dir(entryPath)

	if err := s.store.DeleteEntry(ctx, fname, entryDir); err != nil {
		return fmt.Errorf("delete entry: %w", err)
	}

	return nil
}
