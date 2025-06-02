package userfileFs

import (
	"context"
	"errors"
	"fmt"
	"os"

	"upfile/internal/userfile"
)

type Store struct{}

func New() *Store {
	return &Store{}
}

func (s Store) ReadFile(ctx context.Context, path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", userfile.ErrNotFound
		}

		return "", fmt.Errorf("read file: %w", err)
	}

	return string(content), nil
}
