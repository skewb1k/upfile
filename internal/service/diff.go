package service

import (
	"context"
	"errors"
	"fmt"

	"upfile/internal/store"
)

func (s Service) Diff(
	ctx context.Context,
	fname string,
	entryPath string,
) (string, error) {
	headHash, err := s.store.GetHead(ctx, fname)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return "", ErrNotTracked
		}

		return "", fmt.Errorf("get head: %w", err)
	}

	commit, err := s.store.GetCommitByHash(ctx, fname, headHash)
	if err != nil {
		return "", fmt.Errorf("get commit by hash: %w", err)
	}

	return commit.Content, nil
}
