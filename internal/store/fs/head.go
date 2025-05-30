package storeFs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"upfile/internal/store"
)

const headFname = "HEAD"

func (p Store) SetHead(ctx context.Context, fname string, value string) error {
	fpath := filepath.Join(p.BaseDir, fname, headFname)

	if err := os.MkdirAll(filepath.Dir(fpath), 0o755); err != nil {
		return fmt.Errorf("failed to create dir: %w", err)
	}

	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(value); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	return nil
}

func (p Store) GetHead(ctx context.Context, fname string) (string, error) {
	fpath := filepath.Join(p.BaseDir, fname, headFname)

	data, err := os.ReadFile(fpath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", store.ErrNotFound
		}

		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(data), nil
}
