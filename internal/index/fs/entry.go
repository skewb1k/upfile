package indexFs

import (
	"context"

	"github.com/skewb1k/upfile/internal/index"
)

func (s IndexFsProvider) CreateEntry(
	ctx context.Context,
	fname string,
	entry string,
) error {
	byEntryPath := s.getPathToFilenamesByEntry(entry)

	byDir, err := loadEntrySet(byEntryPath)
	if err != nil {
		return err
	}

	if !byDir.Add(fname) {
		return index.ErrExists
	}

	byNamePath := s.getPathToEntriesByName(fname)

	byName, err := loadEntrySet(byNamePath)
	if err != nil {
		return err
	}

	if !byName.Add(entry) {
		return index.ErrExists
	}

	if err := byName.Save(byNamePath); err != nil {
		return err
	}

	if err := byDir.Save(byEntryPath); err != nil {
		return err
	}

	return nil
}

func (s IndexFsProvider) GetEntriesByFilename(ctx context.Context, fname string) ([]string, error) {
	byNamePath := s.getPathToEntriesByName(fname)

	byname, err := loadEntrySet(byNamePath)
	if err != nil {
		return nil, err
	}

	return byname.ToSlice(), nil
}

func (s IndexFsProvider) DeleteEntry(
	ctx context.Context,
	fname string,
	entry string,
) error {
	byEntryPath := s.getPathToFilenamesByEntry(entry)

	byDir, err := loadEntrySet(byEntryPath)
	if err != nil {
		return err
	}

	if !byDir.Delete(fname) {
		return index.ErrNotFound
	}

	byNamePath := s.getPathToEntriesByName(fname)

	byName, err := loadEntrySet(byNamePath)
	if err != nil {
		return err
	}

	if !byName.Delete(entry) {
		return index.ErrNotFound
	}

	if err := byName.Save(byNamePath); err != nil {
		return err
	}

	if err := byDir.Save(byEntryPath); err != nil {
		return err
	}

	return nil
}

func (s IndexFsProvider) CheckEntry(ctx context.Context, fname string, entry string) (bool, error) {
	byNamePath := s.getPathToEntriesByName(fname)

	byName, err := loadEntrySet(byNamePath)
	if err != nil {
		return false, err
	}

	_, exists := byName[entry]
	return exists, nil
}

func (s IndexFsProvider) GetFilenamesByEntry(ctx context.Context, entry string) ([]string, error) {
	byEntryPath := s.getPathToFilenamesByEntry(entry)

	filenames, err := loadEntrySet(byEntryPath)
	if err != nil {
		return nil, err
	}

	if len(filenames) == 0 {
		return nil, index.ErrNotFound
	}

	return filenames.ToSlice(), nil
}
