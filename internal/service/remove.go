package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/index"
)

func Remove(
	ctx context.Context,
	indexProvider IndexProvider,
	path string,
) error {
	entry, fname := filepath.Dir(path), filepath.Base(path)

	if err := indexProvider.DeleteEntry(ctx, fname, entry); err != nil {
		if errors.Is(err, index.ErrNotFound) {
			return ErrNoEntry
		}

		return fmt.Errorf("delete entry: %w", err)
	}

	return nil
}
