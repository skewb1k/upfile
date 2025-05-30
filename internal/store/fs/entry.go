package storeFs

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
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

func (e *Repr) Add(entry string) {
	if _, exists := e.m[entry]; !exists {
		e.m[entry] = struct{}{}
		e.l = append(e.l, entry)
	}
}

func (e Repr) Save(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("failed to create dir: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	for _, entry := range e.l {
		if _, err := f.WriteString(entry + "\n"); err != nil {
			return fmt.Errorf("failed to write string: %w", err)
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

		return &Repr{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	entries := Repr{
		m: make(map[string]struct{}),
		l: make([]string, 0),
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		entries.Add(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return &Repr{}, fmt.Errorf("failed to write string: %w", err)
	}

	return &entries, nil
}

const entriesFname = "ENTRIES"

func (p Store) CreateEntry(
	ctx context.Context,
	fname string,
	entry string,
) error {
	fpath := filepath.Join(p.BaseDir, fname, entriesFname)
	repr, err := loadRepr(fpath)
	if err != nil {
		return fmt.Errorf("failed to load heads file: %w", err)
	}

	repr.Add(entry)

	if err := repr.Save(fpath); err != nil {
		return fmt.Errorf("failed to write heads file: %w", err)
	}

	return nil
}

func (p Store) CheckEntry(ctx context.Context, fname string, entry string) (bool, error) {
	repr, err := loadRepr(filepath.Join(p.BaseDir, fname, entriesFname))
	if err != nil {
		return false, fmt.Errorf("failed to load heads file: %w", err)
	}

	_, exists := repr.m[entry]
	return exists, nil
}

func (p Store) GetEntries(ctx context.Context, fname string) ([]string, error) {
	repr, err := loadRepr(filepath.Join(p.BaseDir, fname, entriesFname))
	if err != nil {
		return nil, fmt.Errorf("failed to load heads file: %w", err)
	}

	return repr.l, nil
}

func (p Store) GetFiles(ctx context.Context) ([]string, error) {
	entries, err := os.ReadDir(p.BaseDir)
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
