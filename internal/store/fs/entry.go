package storeFs

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"upfile/internal/store"
)

type Repr struct {
	m map[string]struct{}
	l []string
}

func newRepr() *Repr {
	return &Repr{
		m: make(map[string]struct{}),
		l: make([]string, 0),
	}
}

func (r *Repr) Add(entry string) bool {
	if _, exists := r.m[entry]; exists {
		return false
	}

	r.m[entry] = struct{}{}
	r.l = append(r.l, entry)
	return true
}

func (r Repr) Save(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	for _, entry := range r.l {
		if _, err := f.WriteString(entry + "\n"); err != nil {
			return fmt.Errorf("write entry: %w", err)
		}
	}

	return nil
}

func loadRepr(path string) (*Repr, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return newRepr(), nil
		}

		return &Repr{}, fmt.Errorf("open entries file: %w", err)
	}
	defer f.Close()

	entries := Repr{
		m: make(map[string]struct{}),
		l: make([]string, 0),
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		_ = entries.Add(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return &Repr{}, fmt.Errorf("write entry: %w", err)
	}

	return &entries, nil
}

const entriesFname = "ENTRIES"

func (s Store) CreateEntry(
	ctx context.Context,
	fname string,
	entry string,
) error {
	fpath := filepath.Join(s.BaseDir, fname, entriesFname)
	repr, err := loadRepr(fpath)
	if err != nil {
		return err
	}

	if !repr.Add(entry) {
		return store.ErrExists
	}

	if err := repr.Save(fpath); err != nil {
		return err
	}

	return nil
}

func (s Store) CheckEntry(ctx context.Context, fname string, entry string) (bool, error) {
	repr, err := loadRepr(filepath.Join(s.BaseDir, fname, entriesFname))
	if err != nil {
		return false, err
	}

	_, exists := repr.m[entry]
	return exists, nil
}

func (s Store) GetEntries(ctx context.Context, fname string) ([]string, error) {
	repr, err := loadRepr(filepath.Join(s.BaseDir, fname, entriesFname))
	if err != nil {
		return nil, err
	}

	return repr.l, nil
}

func (s Store) GetFiles(ctx context.Context) ([]string, error) {
	entries, err := os.ReadDir(s.BaseDir)
	if err != nil {
		return nil, fmt.Errorf("read base dir: %w", err)
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}

	return dirs, nil
}
