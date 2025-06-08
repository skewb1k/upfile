package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/index"
	"github.com/skewb1k/upfile/internal/userfile"
)

func (s Service) Pull(
	ctx context.Context,
	path string,
) error {
	fname, destDir := filepath.Base(path), filepath.Dir(path)

	upstream, err := s.indexProvider.GetUpstream(ctx, fname)
	if err != nil {
		return fmt.Errorf("get upstream: %w", err)
	}

	existing, err := s.userfileProvider.ReadFile(ctx, path)
	if err != nil {
		if !errors.Is(err, userfile.ErrNotFound) {
			return fmt.Errorf("check existing file: %w", err)
		}
	} else {
		if upstream.Hash.EqualString(existing) {
			return ErrUpToDate
		}
	}

	if err := s.indexProvider.CreateEntry(ctx, fname, destDir); err != nil {
		if !errors.Is(err, index.ErrExists) {
			return fmt.Errorf("create entry: %w", err)
		}
	}

	if err := s.userfileProvider.WriteFile(ctx, path, upstream.Content); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
