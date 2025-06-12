package store

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/skewb1k/upfile/pkg/sha256"
)

type Store struct {
	BaseDir string
}

func New(baseDir string) *Store {
	return &Store{
		BaseDir: baseDir,
	}
}

func (s Store) SaveBlob(ctx context.Context, hash sha256.SHA256, content []byte) error {
	hashStr := hex.EncodeToString(hash[:])
	dir := filepath.Join(s.getPathToBlobs(), hashStr[:2])
	file := filepath.Join(dir, hashStr[2:])

	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	gz := gzip.NewWriter(f)
	defer gz.Close()

	if _, err := gz.Write(content); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

func (p Store) GetBlob(ctx context.Context, hash sha256.SHA256) ([]byte, error) {
	hashStr := hex.EncodeToString(hash[:])
	dir := filepath.Join(p.getPathToBlobs(), hashStr[:2])
	file := filepath.Join(dir, hashStr[2:])

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("gzip new reader: %w", err)
	}
	defer gz.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, gz); err != nil {
		return nil, fmt.Errorf("copy to bytes.Buffer: %w", err)
	}

	return buf.Bytes(), nil
}
