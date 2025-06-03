package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"upfile/internal/index"
	"upfile/internal/userfile"
)

func (s Service) Add(
	ctx context.Context,
	path string,
) error {
	content, err := s.userfileProvider.ReadFile(ctx, path)
	if err != nil {
		if errors.Is(err, userfile.ErrNotFound) {
			return ErrFileNotFound
		}

		return fmt.Errorf("read file: %w", err)
	}

	fname, entryDir := filepath.Base(path), filepath.Dir(path)

	if err := s.indexProvider.CreateEntry(ctx, fname, entryDir); err != nil {
		if errors.Is(err, index.ErrExists) {
			return ErrAlreadyTracked
		}

		return fmt.Errorf("create entry: %w", err)
	}

	upstreamExists, err := s.indexProvider.CheckUpstream(ctx, fname)
	if err != nil {
		return fmt.Errorf("check upstream: %w", err)
	}

	if !upstreamExists {
		if err := s.indexProvider.SetUpstream(ctx, fname, content); err != nil {
			return fmt.Errorf("set upstream: %w", err)
		}
	}

	return nil
}
