package store

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/commit"
	"github.com/skewb1k/upfile/pkg/sha256"
)

func (s Store) SaveCommit(ctx context.Context, hash sha256.SHA256, c *commit.Commit) error {
	hashStr := hash.Hex()
	dir := filepath.Join(s.getPathToCommits(), hashStr[:2])
	file := filepath.Join(dir, hashStr[2:])

	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}

	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	c.EncodeCommit(f)
	return nil
}

func (s Store) GetCommit(ctx context.Context, hash sha256.SHA256) (commit.Commit, error) {
	hashStr := hash.Hex()
	dir := filepath.Join(s.getPathToCommits(), hashStr[:2])
	file := filepath.Join(dir, hashStr[2:])

	f, err := os.Open(file)
	if err != nil {
		return commit.Commit{}, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	c, err := commit.DecodeCommit(f)
	if err != nil {
		return commit.Commit{}, err
	}

	return *c, nil
}
