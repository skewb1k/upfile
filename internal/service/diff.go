package service

import (
	"context"
	"fmt"
)

func (s Service) Diff(
	ctx context.Context,
	fname string,
) (string, error) {
	content, err := s.indexStore.GetUpstream(ctx, fname)
	if err != nil {
		return "", fmt.Errorf("get upstream: %w", err)
	}

	return content, nil
}
