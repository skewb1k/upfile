package indexFs

import (
	"bufio"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/index"
)

type Provider struct {
	BaseDir string
}

func New(baseDir string) *Provider {
	return &Provider{
		BaseDir: baseDir,
	}
}

func (p Provider) GetFilenames(ctx context.Context) ([]string, error) {
	entries, err := os.ReadDir(p.getUpstreams())
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

func (p Provider) SetUpstream(ctx context.Context, fname string, upstream *index.Upstream) error {
	path, err := p.getPathToUpstream(fname)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create versions dir: %w", err)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	gw := gzip.NewWriter(f)
	defer gw.Close()

	if _, err := gw.Write(upstream.Hash[:]); err != nil {
		return fmt.Errorf("write hash: %w", err)
	}

	if _, err := gw.Write([]byte(upstream.Content)); err != nil {
		return fmt.Errorf("gzip write: %w", err)
	}

	return nil
}

func (p Provider) GetUpstream(ctx context.Context, fname string) (index.Upstream, error) {
	path, err := p.getPathToUpstream(fname)
	if err != nil {
		return index.Upstream{}, err
	}

	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return index.Upstream{}, index.ErrNotFound
		}
		return index.Upstream{}, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		return index.Upstream{}, fmt.Errorf("gzip reader: %w", err)
	}
	defer gr.Close()

	r := bufio.NewReader(gr)

	var hash [32]byte
	if _, err := io.ReadFull(r, hash[:]); err != nil {
		return index.Upstream{}, fmt.Errorf("read hash: %w", err)
	}

	content, err := io.ReadAll(r)
	if err != nil {
		return index.Upstream{}, fmt.Errorf("read content: %w", err)
	}

	return index.Upstream{
		Hash:    hash,
		Content: string(content),
	}, nil
}

func (p Provider) CheckUpstream(ctx context.Context, fname string) (bool, error) {
	path, err := p.getPathToUpstream(fname)
	if err != nil {
		return false, err
	}

	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, fmt.Errorf("read file: %w", err)
	}

	return true, nil
}

func (p Provider) DeleteUpstream(ctx context.Context, fname string) error {
	path, err := p.getPathToUpstream(fname)
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return index.ErrNotFound
		}

		return fmt.Errorf("delete file: %w", err)
	}

	return nil
}
