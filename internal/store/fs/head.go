package storeFs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"upfile/internal/store"
)

const headFname = "HEAD"

func (s Store) SetHeadIfNotExists(ctx context.Context, fname string, value string) error {
	headPath := filepath.Join(s.BaseDir, fname, headFname)

	if err := os.MkdirAll(filepath.Dir(headPath), 0o755); err != nil {
		return fmt.Errorf("create head dir: %w", err)
	}

	f, err := os.OpenFile(headPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o644)
	if err != nil {
		if os.IsExist(err) {
			return store.ErrExists
		}
		return fmt.Errorf("create head file: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(value); err != nil {
		return fmt.Errorf("write head: %w", err)
	}

	return nil
}

func (s Store) SetHead(ctx context.Context, fname string, value string) error {
	headPath := filepath.Join(s.BaseDir, fname, headFname)

	if err := os.MkdirAll(filepath.Dir(headPath), 0o755); err != nil {
		return fmt.Errorf("create head dir: %w", err)
	}

	f, err := os.OpenFile(headPath, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("create head file: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(value); err != nil {
		return fmt.Errorf("write head: %w", err)
	}

	return nil
}

func (s Store) GetHead(ctx context.Context, fname string) (string, error) {
	headPath := filepath.Join(s.BaseDir, fname, headFname)

	data, err := os.ReadFile(headPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", store.ErrNotFound
		}

		return "", fmt.Errorf("read head file: %w", err)
	}

	return string(data), nil
}
