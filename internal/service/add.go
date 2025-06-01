package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

func (s Service) Add(
	ctx context.Context,
	path string,
) error {
	fname := filepath.Base(path)
	entryDir := filepath.Dir(path)

	entryExists, err := s.store.CheckEntry(ctx, fname, entryDir)
	if err != nil {
		return fmt.Errorf("get entry: %w", err)
	}

	if entryExists {
		return ErrAlreadyTracked
	}

	if err := s.store.CreateEntry(ctx, fname, entryDir); err != nil {
		return fmt.Errorf("create entry: %w", err)
	}

	upstremExists, err := s.store.CheckUpstream(ctx, fname)
	if err != nil {
		return fmt.Errorf("check upstream: %w", err)
	}

	if !upstremExists {
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read file: %w", err)
		}

		if err := s.store.SetUpstream(ctx, fname, string(content)); err != nil {
			return fmt.Errorf("set upstream: %w", err)
		}
	}

	return nil
}
