package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/skewb1k/upfile/internal/index"
)

func (s Service) Show(
	ctx context.Context,
	fname string,
) (string, error) {
	upstream, err := s.indexProvider.GetUpstream(ctx, fname)
	if err != nil {
		if errors.Is(err, index.ErrNotFound) || errors.Is(err, index.ErrInvalidFilename) {
			return "", ErrNotTracked
		}

		return "", fmt.Errorf("get upstream: %w", err)
	}

	return upstream.Content, nil
}
