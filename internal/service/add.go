package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/index"
)

func Add(
	ctx context.Context,
	indexProvider IndexProvider,
	path string,
) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	entry, fname := filepath.Dir(path), filepath.Base(path)

	if err := indexProvider.CreateEntry(ctx, fname, entry); err != nil {
		if errors.Is(err, index.ErrExists) {
			return ErrAlreadyTracked
		}

		return fmt.Errorf("create entry: %w", err)
	}

	upstreamExists, err := indexProvider.CheckUpstream(ctx, fname)
	if err != nil {
		return fmt.Errorf("check upstream: %w", err)
	}

	if !upstreamExists {
		if err := indexProvider.SetUpstream(ctx, fname, New(content)); err != nil {
			return fmt.Errorf("set upstream: %w", err)
		}
	}

	return nil
}
