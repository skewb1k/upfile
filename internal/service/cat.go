package service

import (
	"context"
	"fmt"
)

func (s Service) CatLatest(ctx context.Context, fname string) ([]byte, error) {
	head, err := s.GetHead(ctx, fname)
	if err != nil {
		return nil, fmt.Errorf("get head: %w", err)
	}

	headCommit, err := s.GetCommit(ctx, head)
	if err != nil {
		return nil, fmt.Errorf("get commit: %w", err)
	}

	content, err := s.GetBlob(ctx, headCommit.ContentHash)
	if err != nil {
		return nil, fmt.Errorf("get blob: %w", err)
	}

	return content, nil
}
