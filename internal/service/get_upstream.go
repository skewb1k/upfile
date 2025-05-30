package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"upfile/internal/store"
)

func GetUpstream(
	ctx context.Context,
	s store.Provider,
	relPath string,
) (string, error) {
	fname := filepath.Base(relPath)
	headHash, err := s.GetHead(ctx, fname)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return "", ErrNotTracked
		}

		return "", fmt.Errorf("get head: %w", err)
	}

	commit, err := s.GetCommitByHash(ctx, fname, headHash)
	if err != nil {
		return "", fmt.Errorf("get commit by hash: %w", err)
	}

	return commit.Content, nil
}
