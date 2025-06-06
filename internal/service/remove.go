package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/index"
)

func (s Service) Remove(
	ctx context.Context,
	path string,
) error {
	fname := filepath.Base(path)
	entryDir := filepath.Dir(path)

	if err := s.indexProvider.DeleteEntry(ctx, fname, entryDir); err != nil {
		if errors.Is(err, index.ErrNotFound) {
			return ErrNotTracked
		}

		return fmt.Errorf("delete entry: %w", err)
	}

	return nil
}
