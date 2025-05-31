package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"upfile/internal/store"
)

func (s Service) Pull(
	ctx context.Context,
	fname string,
	destDir string,
) (bool, error) {
	headHash, err := s.store.GetHead(ctx, fname)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return false, ErrNotTracked
		}

		return false, fmt.Errorf("get head: %w", err)
	}

	commit, err := s.store.GetCommitByHash(ctx, fname, headHash)
	if err != nil {
		return false, fmt.Errorf("get commit by hash: %w", err)
	}

	destPath := filepath.Join(destDir, fname)

	existing, err := os.ReadFile(destPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return false, fmt.Errorf("check existing file: %w", err)
		}
	} else {
		if computeHash(existing) == headHash {
			return false, nil
		}
	}

	if err := s.store.CreateEntry(ctx, fname, destDir); err != nil {
		if !errors.Is(err, store.ErrExists) {
			return false, fmt.Errorf("create entry: %w", err)
		}
	}

	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return false, fmt.Errorf("create parent dirs: %w", err)
	}

	if err := os.WriteFile(destPath, []byte(commit.Content), 0o644); err != nil {
		return false, fmt.Errorf("write file: %w", err)
	}

	return true, nil
}
