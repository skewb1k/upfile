package service

import (
	"context"
	"fmt"

	"upfile/internal/store"
)

func GetFiles(
	ctx context.Context,
	s store.Provider,
) ([]string, error) {
	files, err := s.GetFiles(ctx)
	if err != nil {
		return nil, fmt.Errorf("get files: %w", err)
	}

	return files, nil
}
