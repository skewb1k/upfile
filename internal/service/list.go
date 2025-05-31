package service

import (
	"context"
	"fmt"
)

func (s Service) List(
	ctx context.Context,
) ([]string, error) {
	files, err := s.store.GetFiles(ctx)
	if err != nil {
		return nil, fmt.Errorf("get files: %w", err)
	}

	return files, nil
}
