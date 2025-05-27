package fs

import (
	"context"
	"fmt"
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

const entriesFname = "ENTRIES"

func (s Store) CreateEntry(
	ctx context.Context,
	fname string,
	entry string,
) error {
	fpath := filepath.Join(s.BaseDir, fname, entriesFname)
	e, err := Load(fpath)
	if err != nil {
		return fmt.Errorf("failed to load heads file: %w", err)
	}

	e.Add(entry)

	if err := e.Save(fpath); err != nil {
		return fmt.Errorf("failed to write heads file: %w", err)
	}

	return nil
}

func (s Store) CheckEntry(ctx context.Context, fname string, entry string) (bool, error) {
	e, err := Load(filepath.Join(s.BaseDir, fname, entriesFname))
	if err != nil {
		return false, fmt.Errorf("failed to load heads file: %w", err)
	}

	_, exists := e.m[entry]
	return exists, nil
}

func (s Store) GetEntries(ctx context.Context, fname string) ([]string, error) {
	e, err := Load(filepath.Join(s.BaseDir, fname, entriesFname))
	if err != nil {
		return nil, fmt.Errorf("failed to load heads file: %w", err)
	}

	return e.l, nil
}
