package entries

import (
	"context"
)

type Provider struct {
	BaseDir string
}

func NewProvider(baseDir string) *Provider {
	return &Provider{
		BaseDir: baseDir,
	}
}

func (p Provider) CreateEntry(
	ctx context.Context,
	fname string,
	entry string,
) error {
	byEntryPath := p.getPathToFilenamesByEntry(entry)

	bydir, err := Load(byEntryPath)
	if err != nil {
		return err
	}

	if !bydir.Add(fname) {
		return ErrExists
	}

	byNamePath := p.getPathToEntriesByName(fname)

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
	byEntryPath := p.getPathToFilenamesByEntry(entry)

	bydir, err := Load(byEntryPath)
	if err != nil {
		return err
	}

	if !bydir.Delete(fname) {
		return ErrNotFound
	}

	byNamePath := p.getPathToEntriesByName(fname)

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

func (p Provider) CheckEntry(ctx context.Context, fname string, entry string) (bool, error) {
	byNamePath := p.getPathToEntriesByName(fname)

	byname, err := Load(byNamePath)
	if err != nil {
		return false, err
	}

	_, exists := byname[entry]
	return exists, nil
}

func (p Provider) GetFilenamesByEntry(ctx context.Context, entry string) ([]string, error) {
	byEntryPath := p.getPathToFilenamesByEntry(entry)

	filenames, err := Load(byEntryPath)
	if err != nil {
		return nil, err
	}

	if len(filenames) == 0 {
		return nil, ErrNotFound
	}

	return filenames.ToSlice(), nil
}
