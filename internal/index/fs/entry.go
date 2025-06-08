package indexFs

import (
	"context"

	"github.com/skewb1k/upfile/internal/index"
)

func (p Provider) CreateEntry(
	ctx context.Context,
	fname string,
	entry string,
) error {
	byEntryPath := p.getPathToFilenamesByEntry(encodePath(entry))

	bydir, err := Load(byEntryPath)
	if err != nil {
		return err
	}

	if !bydir.Add(fname) {
		return index.ErrExists
	}

	byNamePath, err := p.getPathToEntriesByName(fname)
	if err != nil {
		return err
	}

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
	byNamePath, err := p.getPathToEntriesByName(fname)
	if err != nil {
		return nil, err
	}

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
	byEntryPath := p.getPathToFilenamesByEntry(encodePath(entry))

	bydir, err := Load(byEntryPath)
	if err != nil {
		return err
	}

	if !bydir.Delete(fname) {
		return index.ErrNotFound
	}

	byNamePath, err := p.getPathToEntriesByName(fname)
	if err != nil {
		return err
	}

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
	byNamePath, err := p.getPathToEntriesByName(fname)
	if err != nil {
		return false, err
	}

	byname, err := Load(byNamePath)
	if err != nil {
		return false, err
	}

	_, exists := byname[entry]
	return exists, nil
}

func (p Provider) GetFilenamesByEntry(ctx context.Context, entry string) ([]string, error) {
	byEntryPath := p.getPathToFilenamesByEntry(encodePath(entry))

	filenames, err := Load(byEntryPath)
	if err != nil {
		return nil, err
	}

	if len(filenames) == 0 {
		return nil, index.ErrNotFound
	}

	return filenames.ToSlice(), nil
}
