package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

func (s Service) Push(
	ctx context.Context,
	path string,
) error {
	fname := filepath.Base(path)
	// entryDir := filepath.Dir(path)

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	current, err := s.store.GetUpstream(ctx, fname)
	if err != nil {
		return fmt.Errorf("get upstream: %w", err)
	}

	if slices.Equal(content, []byte(current)) {
		return ErrUpToDate
	}

	if err := s.store.SetUpstream(ctx, fname, string(content)); err != nil {
		return fmt.Errorf("set upstream: %w", err)
	}

	return nil
}
