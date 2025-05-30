package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"upfile/internal/store"
)

func Link(
	ctx context.Context,
	s store.Provider,
	path string,
) error {
	fname := filepath.Base(path)
	entryDir := filepath.Dir(path)

	entryExists, err := s.CheckEntry(ctx, fname, entryDir)
	if err != nil {
		return fmt.Errorf("failed to get entry: %w", err)
	}

	if entryExists {
		return ErrAlreadyLinked
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	hash := computeHash(content)

	if _, err := s.GetCommitByHash(ctx, fname, hash); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			if err := s.CreateCommit(ctx, fname, &store.Commit{
				Hash:    hash,
				Content: string(content),
				Parent:  "",
			}); err != nil {
				return fmt.Errorf("failed to create commit: %w", err)
			}
		} else {
			return fmt.Errorf("failed to get commit: %w", err)
		}
	}

	if err := s.CreateEntry(ctx, fname, entryDir); err != nil {
		return fmt.Errorf("failed to create entry: %w", err)
	}

	if err := s.SetHead(ctx, fname, hash); err != nil {
		return fmt.Errorf("failed to set head: %w", err)
	}

	return nil
}
