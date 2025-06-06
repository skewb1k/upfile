package userfileFs

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/userfile"
)

type Provider struct{}

func New() *Provider {
	return &Provider{}
}

func (p Provider) ReadFile(ctx context.Context, path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", userfile.ErrNotFound
		}

		return "", fmt.Errorf("read file: %w", err)
	}

	return string(content), nil
}

func (p Provider) WriteFile(ctx context.Context, path string, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create parent dirs: %w", err)
	}

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

func (p Provider) CheckFile(ctx context.Context, path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, fmt.Errorf("open file: %w", err)
	}

	_ = f.Close()

	return true, nil
}
