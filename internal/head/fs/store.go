package fs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

type Store struct {
	BaseDir string
}

func NewStore(baseDir string) *Store {
	return &Store{
		BaseDir: baseDir,
	}
}

const headFname = "HEAD"

func (s Store) SetHead(ctx context.Context, fname string, value string) error {
	fpath := filepath.Join(s.BaseDir, fname, headFname)

	if err := os.MkdirAll(filepath.Dir(fpath), 0o700); err != nil {
		return fmt.Errorf("failed to create dir: %w", err)
	}

	f, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(value); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}

	return nil
}

func (s Store) GetHead(ctx context.Context, fname string) (string, error) {
	fpath := filepath.Join(s.BaseDir, fname, headFname)

	data, err := os.ReadFile(fpath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}

		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(data), nil
}
