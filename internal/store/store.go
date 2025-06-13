package store

import (
	"context"
)

type Store struct {
	BaseDir string
}

func New(baseDir string) *Store {
	return &Store{
		BaseDir: baseDir,
	}
}

func (s Store) CreateEntry(
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
		return ErrExists
	}

	byNamePath := s.getPathToEntriesByName(fname)

	byName, err := loadEntrySet(byNamePath)
	if err != nil {
		return err
	}

	if !byName.Add(entry) {
		return ErrExists
	}

	if err := byName.Save(byNamePath); err != nil {
		return err
	}

	if err := byDir.Save(byEntryPath); err != nil {
		return err
	}

	return nil
}

func (s Store) GetEntriesByFilename(ctx context.Context, fname string) ([]string, error) {
	byNamePath := s.getPathToEntriesByName(fname)

	byname, err := loadEntrySet(byNamePath)
	if err != nil {
		return nil, err
	}

	return byname.ToSlice(), nil
}

func (s Store) DeleteEntry(
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
		return ErrNotFound
	}

	byNamePath := s.getPathToEntriesByName(fname)

	byName, err := loadEntrySet(byNamePath)
	if err != nil {
		return err
	}

	if !byName.Delete(entry) {
		return ErrNotFound
	}

	if err := byName.Save(byNamePath); err != nil {
		return err
	}

	if err := byDir.Save(byEntryPath); err != nil {
		return err
	}

	return nil
}

func (s Store) CheckEntry(ctx context.Context, fname string, entry string) (bool, error) {
	byNamePath := s.getPathToEntriesByName(fname)

	byName, err := loadEntrySet(byNamePath)
	if err != nil {
		return false, err
	}

	_, exists := byName[entry]
	return exists, nil
}

func (s Store) GetFilenamesByEntry(ctx context.Context, entry string) ([]string, error) {
	byEntryPath := s.getPathToFilenamesByEntry(entry)

	filenames, err := loadEntrySet(byEntryPath)
	if err != nil {
		return nil, err
	}

	if len(filenames) == 0 {
		return nil, ErrNotFound
	}

	return filenames.ToSlice(), nil
}
