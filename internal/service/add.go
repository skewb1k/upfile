package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"upfile/internal/store"
)

func (s Service) Add(
	ctx context.Context,
	entryPath string,
) error {
	fname := filepath.Base(entryPath)
	entryDir := filepath.Dir(entryPath)

	entryExists, err := s.store.CheckEntry(ctx, fname, entryDir)
	if err != nil {
		return fmt.Errorf("get entry: %w", err)
	}

	if entryExists {
		return ErrAlreadyLinked
	}

	content, err := os.ReadFile(entryPath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	hash := computeHash(content)

	if _, err := s.store.GetCommitByHash(ctx, fname, hash); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			if err := s.store.CreateCommit(ctx, fname, &store.Commit{
				Hash:    hash,
				Content: string(content),
				Parent:  "",
			}); err != nil {
				return fmt.Errorf("create commit: %w", err)
			}
		} else {
			return fmt.Errorf("get commit by hash: %w", err)
		}
	}

	if err := s.store.CreateEntry(ctx, fname, entryDir); err != nil {
		return fmt.Errorf("create entry: %w", err)
	}

	if err := s.store.SetHeadIfNotExists(ctx, fname, hash); err != nil {
		if !errors.Is(err, store.ErrExists) {
			return fmt.Errorf("set head: %w", err)
		}
	}

	return nil
}
