package indexFs

import (
	"context"
	"encoding/base64"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/index"
)

const (
	entriesDirname = "entries"
	byEntry        = "by-entry"
	byName         = "by-name"
)

func (p Provider) getPathToEntriesByName(fname string) string {
	return filepath.Join(
		p.BaseDir,
		entriesDirname,
		byName,
		filepath.Clean(fname),
	)
}

func (p Provider) getPathToFilenamesByEntries(entry string) string {
	return filepath.Join(
		p.BaseDir,
		entriesDirname,
		byEntry,
		filepath.Clean(entry),
	)
}

func encodePath(path string) string {
	return base64.URLEncoding.EncodeToString([]byte(path))
}

//	func decodePath(encoded string) (string, error) {
//		res, err := base64.URLEncoding.DecodeString(encoded)
//		if err != nil {
//			return "", fmt.Errorf("failed to decode base64: %w", err)
//		}
//
//		return string(res), nil
//	}

func (p Provider) CreateEntry(
	ctx context.Context,
	fname string,
	entry string,
) error {
	encoded := encodePath(entry)

	byEntryPath := p.getPathToFilenamesByEntries(encoded)

	bydir, err := Load(byEntryPath)
	if err != nil {
		return err
	}

	if !bydir.Add(fname) {
		return index.ErrExists
	}

	byNamePath := p.getPathToEntriesByName(fname)

	byname, err := Load(byNamePath)
	if err != nil {
		return err
	}

	if !byname.Add(entry) {
		return index.ErrExists
	}

	if err := byname.Save(byNamePath); err != nil {
		return err
	}

	if err := bydir.Save(byEntryPath); err != nil {
		return err
	}

	return nil
}

func (p Provider) GetEntriesByFilename(ctx context.Context, fname string) ([]string, error) {
	byNamePath := p.getPathToEntriesByName(fname)

	byname, err := Load(byNamePath)
	if err != nil {
		return nil, err
	}

	return byname.ToSlice(), nil
}

func (p Provider) DeleteEntry(
	ctx context.Context,
	fname string,
	entry string,
) error {
	encoded := encodePath(entry)

	byEntryPath := p.getPathToFilenamesByEntries(encoded)

	bydir, err := Load(byEntryPath)
	if err != nil {
		return err
	}

	if !bydir.Delete(fname) {
		return index.ErrNotFound
	}

	byNamePath := p.getPathToEntriesByName(fname)

	byname, err := Load(byNamePath)
	if err != nil {
		return err
	}

	if !byname.Delete(entry) {
		return index.ErrNotFound
	}

	if err := byname.Save(byNamePath); err != nil {
		return err
	}

	if err := bydir.Save(byEntryPath); err != nil {
		return err
	}

	return nil
}

func (p Provider) CheckEntry(ctx context.Context, fname string, entry string) (bool, error) {
	byname, err := Load(p.getPathToEntriesByName(fname))
	if err != nil {
		return false, err
	}

	_, exists := byname[entry]
	return exists, nil
}

func (p Provider) GetFilenamesByEntry(ctx context.Context, entry string) ([]string, error) {
	encoded := encodePath(entry)
	byEntryPath := p.getPathToFilenamesByEntries(encoded)

	bydir, err := Load(byEntryPath)
	if err != nil {
		return nil, err
	}

	return bydir.ToSlice(), nil
}
