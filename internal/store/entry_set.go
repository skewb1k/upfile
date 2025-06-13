package store

import (
	"bufio"
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
)

type EntrySet map[string]struct{}

func (e EntrySet) Add(s string) bool {
	if _, exists := e[s]; exists {
		return false
	}

	e[s] = struct{}{}
	return true
}

func (e EntrySet) Delete(s string) bool {
	if _, exists := e[s]; !exists {
		return false
	}

	delete(e, s)
	return true
}

// TODO: improve performance
func (e EntrySet) ToSlice() []string {
	return slices.Sorted(maps.Keys(e))
}

func (e EntrySet) Save(path string) error {
	if len(e) == 0 {
		if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("remove file: %w", err)
		}

		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	for _, item := range e.ToSlice() {
		if _, err := f.WriteString(item + "\n"); err != nil {
			return fmt.Errorf("write string: %w", err)
		}
	}

	return nil
}

func loadEntrySet(path string) (EntrySet, error) {
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return make(EntrySet), nil
		}

		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	set := make(EntrySet)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			set[line] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner: %w", err)
	}

	return set, nil
}
