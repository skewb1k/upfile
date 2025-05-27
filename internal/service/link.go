package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"upfile/internal/commits"
)

func (s Service) Link(ctx context.Context, path string) error {
	fname := filepath.Base(path)
	entryDir := filepath.Dir(path)

	exists, err := s.entries.CheckEntry(ctx, fname, entryDir)
	if err != nil {
		return fmt.Errorf("failed to get entry: %w", err)
	}

	if exists {
		return ErrAlreadyLinked
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	hash := computeHash(data)

	if _, err := s.commits.GetCommitByHash(ctx, fname, hash); err != nil {
		if errors.Is(err, commits.ErrNotFound) {
			if err := s.commits.CreateCommit(ctx, fname, &commits.Commit{
				Hash:    hash,
				Content: string(data),
				Parent:  "",
			}); err != nil {
				return fmt.Errorf("failed to create commit: %w", err)
			}
		} else {
			return fmt.Errorf("failed to get commit: %w", err)
		}
	}

	if err := s.entries.CreateEntry(ctx, fname, entryDir); err != nil {
		return fmt.Errorf("failed to create entry: %w", err)
	}

	if err := s.head.SetHead(ctx, fname, hash); err != nil {
		return fmt.Errorf("failed to set head: %w", err)
	}

	return nil
}
