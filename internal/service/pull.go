package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"upfile/internal/index"
	"upfile/internal/userfile"
)

func (s Service) Pull(
	ctx context.Context,
	path string,
) error {
	fname, destDir := filepath.Base(path), filepath.Dir(path)

	content, err := s.indexProvider.GetUpstream(ctx, fname)
	if err != nil {
		return fmt.Errorf("get upstream: %w", err)
	}

	existing, err := s.userfileProvider.ReadFile(ctx, path)
	if err != nil {
		if !errors.Is(err, userfile.ErrNotFound) {
			return fmt.Errorf("check existing file: %w", err)
		}
	} else {
		if hash(content) == hash(existing) {
			return ErrUpToDate
		}
	}

	if err := s.indexProvider.CreateEntry(ctx, fname, destDir); err != nil {
		if !errors.Is(err, index.ErrExists) {
			return fmt.Errorf("create entry: %w", err)
		}
	}

	if err := s.userfileProvider.WriteFile(ctx, path, content); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
