package service

import (
	"context"
	"errors"
	"fmt"

	"upfile/internal/index"
)

func (s Service) Show(
	ctx context.Context,
	fname string,
) (string, error) {
	content, err := s.indexProvider.GetUpstream(ctx, fname)
	if err != nil {
		if errors.Is(err, index.ErrNotFound) {
			return "", ErrNotTracked
		}
		return "", fmt.Errorf("get upstream: %w", err)
	}

	return content, nil
}
