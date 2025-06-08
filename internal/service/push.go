package service

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/index"
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

	newContent, err := s.userfileProvider.ReadFile(ctx, path)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	upstream, err := s.indexProvider.GetUpstream(ctx, fname)
	if err != nil {
		return fmt.Errorf("get upstream: %w", err)
	}

	if upstream.Hash.EqualString(newContent) {
		return ErrUpToDate
	}

	if err := s.indexProvider.SetUpstream(ctx, fname, index.NewUpstream(newContent)); err != nil {
		return fmt.Errorf("set upstream: %w", err)
	}

	return nil
}
