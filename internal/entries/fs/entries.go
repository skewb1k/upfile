package fs

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

type Repr struct {
	m map[string]struct{}
	l []string
}

func NewRepr() *Repr {
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

func Load(path string) (*Repr, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return NewRepr(), nil
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

func (e Repr) Save(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
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
