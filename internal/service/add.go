package service

import (
	"context"
	"fmt"
	"path/filepath"
)

func (s Service) Add(
	ctx context.Context,
	path string,
) error {
	fname, entryDir := filepath.Base(path), filepath.Dir(path)

	entryExists, err := s.indexStore.CheckEntry(ctx, fname, entryDir)
	if err != nil {
		return fmt.Errorf("get entry: %w", err)
	}

	if entryExists {
		return ErrAlreadyTracked
	}

	if err := s.indexStore.CreateEntry(ctx, fname, entryDir); err != nil {
		return fmt.Errorf("create entry: %w", err)
	}

	upstreamExists, err := s.indexStore.CheckUpstream(ctx, fname)
	if err != nil {
		return fmt.Errorf("check upstream: %w", err)
	}

	if !upstreamExists {
		content, err := s.userfileStore.ReadFile(ctx, path)
		if err != nil {
			// TODO: handle not found error
			return fmt.Errorf("read file: %w", err)
		}

		if err := s.indexStore.SetUpstream(ctx, fname, content); err != nil {
			return fmt.Errorf("set upstream: %w", err)
		}
	}

	return nil
}
