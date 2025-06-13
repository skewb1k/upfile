package store

import (
	"bufio"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/skewb1k/upfile/pkg/sha256"
)

type Upstream struct {
	Hash    sha256.SHA256
	Content string
}

func NewUpstream(content string) *Upstream {
	return &Upstream{
		Hash:    sha256.FromString(content),
		Content: content,
	}
}

func (s Store) GetFilenames(ctx context.Context) ([]string, error) {
	entries, err := os.ReadDir(s.getUpstreams())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}

		return nil, fmt.Errorf("read versions dir: %w", err)
	}

	dirs := make([]string, len(entries))
	for i, entry := range entries {
		dirs[i], err = decodePath(entry.Name())
		if err != nil {
			return nil, err
		}
	}

	return dirs, nil
}

func (s Store) SetUpstream(ctx context.Context, fname string, upstream *Upstream) error {
	path := s.getPathToUpstream(fname)

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

func (s Store) GetUpstream(ctx context.Context, fname string) (Upstream, error) {
	path := s.getPathToUpstream(fname)

	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Upstream{}, ErrNotFound
		}

		return Upstream{}, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		return Upstream{}, fmt.Errorf("gzip reader: %w", err)
	}
	defer gr.Close()

	r := bufio.NewReader(gr)

	var hash [32]byte
	if _, err := io.ReadFull(r, hash[:]); err != nil {
		return Upstream{}, fmt.Errorf("read hash: %w", err)
	}

	content, err := io.ReadAll(r)
	if err != nil {
		return Upstream{}, fmt.Errorf("read content: %w", err)
	}

	return Upstream{
		Hash:    hash,
		Content: string(content),
	}, nil
}

func (s Store) CheckUpstream(ctx context.Context, fname string) (bool, error) {
	path := s.getPathToUpstream(fname)

	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, fmt.Errorf("read file: %w", err)
	}

	return true, nil
}

func (s Store) DeleteUpstream(ctx context.Context, fname string) error {
	path := s.getPathToUpstream(fname)

	if err := os.Remove(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNotFound
		}

		return fmt.Errorf("delete file: %w", err)
	}

	return nil
}
