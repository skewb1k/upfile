package indexFs

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"upfile/internal/index"
)

type Provider struct {
	BaseDir string
}

func New(baseDir string) *Provider {
	return &Provider{
		BaseDir: baseDir,
	}
}

const versionsDirname = "versions"

func (p Provider) GetFiles(ctx context.Context) ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(p.BaseDir, versionsDirname))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}

		return nil, fmt.Errorf("read versions dir: %w", err)
	}

	dirs := make([]string, len(entries))
	for i, entry := range entries {
		dirs[i] = entry.Name()
	}

	return dirs, nil
}

func (p Provider) SetUpstream(ctx context.Context, fname string, value string) error {
	versionsDir := filepath.Join(p.BaseDir, versionsDirname)
	if err := os.MkdirAll(versionsDir, 0o755); err != nil {
		return fmt.Errorf("create versions dir: %w", err)
	}

	if err := os.WriteFile(filepath.Join(versionsDir, fname), []byte(value), 0o644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

func (p Provider) GetUpstream(ctx context.Context, fname string) (string, error) {
	data, err := os.ReadFile(filepath.Join(p.BaseDir, versionsDirname, fname))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", index.ErrNotFound
		}

		return "", fmt.Errorf("read file: %w", err)
	}

	return string(data), nil
}

func (p Provider) CheckUpstream(ctx context.Context, fname string) (bool, error) {
	if _, err := os.Stat(filepath.Join(p.BaseDir, versionsDirname, fname)); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, fmt.Errorf("read file: %w", err)
	}

	return true, nil
}
