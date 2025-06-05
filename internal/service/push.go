package service

import (
	"context"
	"fmt"
	"path/filepath"
)

func (s Service) Push(ctx context.Context, path string) error {
	fname, entryDir := filepath.Base(path), filepath.Dir(path)

	exists, err := s.indexProvider.CheckEntry(ctx, fname, entryDir)
	if err != nil {
		return fmt.Errorf("check entry: %w", err)
	}

	if !exists {
		return ErrNotTracked
	}

	content, err := s.userfileProvider.ReadFile(ctx, path)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	current, err := s.indexProvider.GetUpstream(ctx, fname)
	if err != nil {
		return fmt.Errorf("get upstream: %w", err)
	}

	if hash(content) == hash(current) {
		return ErrUpToDate
	}

	if err := s.indexProvider.SetUpstream(ctx, fname, content); err != nil {
		return fmt.Errorf("set upstream: %w", err)
	}

	return nil
}
