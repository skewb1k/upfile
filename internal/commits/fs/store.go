package fs

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"upfile/internal/commits"
)

type Store struct {
	BaseDir string
}

func NewStore(baseDir string) *Store {
	return &Store{
		BaseDir: baseDir,
	}
}

type Commit struct {
	Content string `json:"content"`
	Parent  string `json:"parent"`
}

func (c Commit) Zip() ([]byte, error) {
	data, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer

	gw := gzip.NewWriter(&buf)
	defer gw.Close()

	if _, err := gw.Write(data); err != nil {
		return nil, fmt.Errorf("failed to gzip commit: %w", err)
	}

	return buf.Bytes(), nil
}

func Unzip(r io.Reader) (*Commit, error) {
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gr.Close()

	decoded, err := io.ReadAll(gr)
	if err != nil {
		return nil, fmt.Errorf("failed to read compressed data: %w", err)
	}

	var c Commit
	if err := json.Unmarshal(decoded, &c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal commit: %w", err)
	}

	return &c, nil
}

const commitsDirname = "commits"

func (s Store) CreateCommit(
	ctx context.Context,
	fname string,
	commit *commits.Commit,
) error {
	commitsDir := filepath.Join(s.BaseDir, fname, commitsDirname)

	commitFname := commit.Hash[2:]
	commitDirname := filepath.Join(commitsDir, commit.Hash[:2])

	if err := os.MkdirAll(commitDirname, 0o700); err != nil {
		return fmt.Errorf("failed to create object subdir: %w", err)
	}

	f, err := os.OpenFile(filepath.Join(commitDirname, commitFname), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("file already exists: %w", err)
		}

		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	b, err := Commit{
		Content: commit.Content,
		Parent:  commit.Parent,
	}.Zip()
	if err != nil {
		return fmt.Errorf("failed to compress commit: %w", err)
	}

	if _, err := f.Write(b); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func (s Store) GetCommitByHash(ctx context.Context, fname string, hash string) (commits.Commit, error) {
	commitsDir := filepath.Join(s.BaseDir, fname, commitsDirname)

	commitFname := hash[2:]
	commitDirname := filepath.Join(commitsDir, hash[:2])

	data, err := os.Open(filepath.Join(commitDirname, commitFname))
	if err != nil {
		if os.IsNotExist(err) {
			return commits.Commit{}, commits.ErrNotFound
		}

		return commits.Commit{}, fmt.Errorf("failed to read commit file: %w", err)
	}

	commit, err := Unzip(data)
	if err != nil {
		return commits.Commit{}, fmt.Errorf("failed to parse commit: %w", err)
	}

	return commits.Commit{
		Hash:    hash,
		Content: commit.Content,
		Parent:  commit.Parent,
	}, nil
}
