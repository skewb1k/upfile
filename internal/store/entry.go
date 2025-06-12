package store

import (
	"context"
)

func (s Store) CreateEntry(
	ctx context.Context,
	fname string,
	entry string,
) error {
	byEntryPath := s.getPathToFilenamesByEntry(entry)

	bydir, err := Load(byEntryPath)
	if err != nil {
		return err
	}

	if !bydir.Add(fname) {
		return ErrExists
	}

	byNamePath := s.getPathToEntriesByName(fname)

	byname, err := Load(byNamePath)
	if err != nil {
		return err
	}

	if !byname.Add(entry) {
		return ErrExists
	}

	if err := byname.Save(byNamePath); err != nil {
		return err
	}

	if err := bydir.Save(byEntryPath); err != nil {
		return err
	}

	return nil
}

func (s Store) GetEntriesByFilename(ctx context.Context, fname string) ([]string, error) {
	byNamePath := s.getPathToEntriesByName(fname)

	byname, err := Load(byNamePath)
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

	bydir, err := Load(byEntryPath)
	if err != nil {
		return err
	}

	if !bydir.Delete(fname) {
		return ErrNotFound
	}

	byNamePath := s.getPathToEntriesByName(fname)

	byname, err := Load(byNamePath)
	if err != nil {
		return err
	}

	if !byname.Delete(entry) {
		return ErrNotFound
	}

	if err := byname.Save(byNamePath); err != nil {
		return err
	}

	if err := bydir.Save(byEntryPath); err != nil {
		return err
	}

	return nil
}

func (s Store) CheckEntry(ctx context.Context, fname string, entry string) (bool, error) {
	byNamePath := s.getPathToEntriesByName(fname)

	byname, err := Load(byNamePath)
	if err != nil {
		return false, err
	}

	_, exists := byname[entry]
	return exists, nil
}

func (s Store) GetFilenamesByEntry(ctx context.Context, entry string) ([]string, error) {
	byEntryPath := s.getPathToFilenamesByEntry(entry)

	filenames, err := Load(byEntryPath)
	if err != nil {
		return nil, err
	}

	if len(filenames) == 0 {
		return nil, ErrNotFound
	}

	return filenames.ToSlice(), nil
}
