package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"upfile/internal/index"
)

func (s Service) Pull(
	ctx context.Context,
	fname string,
	destDir string,
) (bool, error) {
	content, err := s.indexStore.GetUpstream(ctx, fname)
	if err != nil {
		return false, fmt.Errorf("get upstream: %w", err)
	}

	destPath := filepath.Join(destDir, fname)

	existing, err := os.ReadFile(destPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return false, fmt.Errorf("check existing file: %w", err)
		}
	} else {
		if slices.Equal(existing, []byte(content)) {
			return false, nil
		}
	}

	if err := s.indexStore.CreateEntry(ctx, fname, destDir); err != nil {
		if !errors.Is(err, index.ErrExists) {
			return false, fmt.Errorf("create entry: %w", err)
		}
	}

	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return false, fmt.Errorf("create parent dirs: %w", err)
	}

	if err := os.WriteFile(destPath, []byte(content), 0o644); err != nil {
		return false, fmt.Errorf("write file: %w", err)
	}

	return true, nil
}
