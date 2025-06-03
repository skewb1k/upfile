package indexFs

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"upfile/internal/index"
)

const (
	entriesDirname = "entries"
	byDirName      = "by-dir"
	byName         = "by-name"
)

func (p Provider) getPathToEntriesByFname(fname string) string {
	return filepath.Join(
		p.BaseDir,
		entriesDirname,
		byName,
		fname,
	)
}

func (p Provider) getPathToEntriesByDir(dir string) string {
	return filepath.Join(
		p.BaseDir,
		entriesDirname,
		byDirName,
		dir,
	)
}

func encodePath(path []byte) string {
	return base64.URLEncoding.EncodeToString(path)
}

func decodePath(encoded string) ([]byte, error) {
	res, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	return res, nil
}

func (p Provider) CreateEntry(
	ctx context.Context,
	fname string,
	entryDir string,
) error {
	encoded := encodePath([]byte(entryDir))

	byDirPath := p.getPathToEntriesByDir(encoded)
	if err := os.MkdirAll(byDirPath, 0o700); err != nil {
		return fmt.Errorf("mkdir by-dir path: %w", err)
	}

	byDirFile, err := os.OpenFile(filepath.Join(byDirPath, fname), os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return index.ErrExists
		}

		return fmt.Errorf("create by-dir entry: %w", err)
	}
	_ = byDirFile.Close()

	byNameDir := p.getPathToEntriesByFname(fname)

	if err := os.MkdirAll(byNameDir, 0o700); err != nil {
		return fmt.Errorf("mkdir by-name dir: %w", err)
	}

	byNameFile, err := os.OpenFile(filepath.Join(byNameDir, encoded), os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return index.ErrExists
		}

		return fmt.Errorf("create by-name entry: %w", err)
	}
	_ = byNameFile.Close()

	return nil
}

func (p Provider) GetEntriesByFname(ctx context.Context, fname string) ([]string, error) {
	entries, err := os.ReadDir(p.getPathToEntriesByFname(fname))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}

		return nil, fmt.Errorf("read by-name path: %w", err)
	}

	result := make([]string, len(entries))
	for i, entry := range entries {
		e, err := decodePath(entry.Name())
		if err != nil {
			return nil, err
		}

		result[i] = string(e)
	}

	return result, nil
}

func (p Provider) DeleteEntry(
	ctx context.Context,
	fname string,
	entryDir string,
) error {
	encoded := encodePath([]byte(entryDir))
	byDirDir := p.getPathToEntriesByDir(encoded)

	if err := os.Remove(filepath.Join(byDirDir, fname)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = index.ErrNotFound
		}

		return fmt.Errorf("remove by-dir entry file: %w", err)
	}

	if entries, err := os.ReadDir(byDirDir); err == nil && len(entries) == 0 {
		if err := os.Remove(byDirDir); err != nil {
			return fmt.Errorf("remove by-dir entry dir: %w", err)
		}
	}

	byNameDir := p.getPathToEntriesByFname(fname)
	if err := os.Remove(filepath.Join(byNameDir, encoded)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = index.ErrNotFound
		}

		return fmt.Errorf("remove by-name entry file: %w", err)
	}

	if entries, err := os.ReadDir(byNameDir); err == nil && len(entries) == 0 {
		if err := os.Remove(byNameDir); err != nil {
			return fmt.Errorf("remove by-name entry dir: %w", err)
		}
	}

	return nil
}

func (p Provider) CheckEntry(ctx context.Context, fname string, entryDir string) (bool, error) {
	byNamePath := filepath.Join(p.getPathToEntriesByFname(fname), encodePath([]byte(entryDir)))

	if _, err := os.Stat(byNamePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, fmt.Errorf("stat by-name file: %w", err)
	}

	return true, nil
}

func (p Provider) GetFilesByEntryDir(ctx context.Context, entryDir string) ([]string, error) {
	encoded := encodePath([]byte(entryDir))
	byDirDir := p.getPathToEntriesByDir(encoded)

	entries, err := os.ReadDir(byDirDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}

		return nil, fmt.Errorf("read versions dir: %w", err)
	}

	dirs := make([]string, len(entries))
	for i, entry := range entries {
		dirs[i] = entry.Name()
	}

	return dirs, nil
}
