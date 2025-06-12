package service

import (
	"context"

	"github.com/skewb1k/upfile/internal/commit"
	"github.com/skewb1k/upfile/pkg/sha256"
)

type IndexProvider interface {
	CreateEntry(ctx context.Context, fname string, entry string) error
	CheckEntry(ctx context.Context, fname string, entry string) (bool, error)
	GetEntriesByFilename(ctx context.Context, fname string) ([]string, error)
	GetFilenamesByEntry(ctx context.Context, entry string) ([]string, error)
	// DeleteEntry(ctx context.Context, fname string, entry string) error
}

type BlobProvider interface {
	SaveBlob(ctx context.Context, hash sha256.SHA256, content []byte) error
	GetBlob(ctx context.Context, hash sha256.SHA256) ([]byte, error)
	// Load(ctx context.Context, hash sha256.SHA256) ([]byte, error)
}

type CommitProvider interface {
	SaveCommit(ctx context.Context, hash sha256.SHA256, commit *commit.Commit) error
	GetCommit(ctx context.Context, hash sha256.SHA256) (commit.Commit, error)
	// ListCommits(filename string) ([]string, error)
}

type HeadProvider interface {
	GetHead(ctx context.Context, filename string) (sha256.SHA256, error)
	CheckHead(ctx context.Context, filename string) (bool, error)
	SetHead(ctx context.Context, filename string, commitHash sha256.SHA256) error
	GetFilenames(ctx context.Context) ([]string, error)
}

type Store interface {
	IndexProvider
	BlobProvider
	CommitProvider
	HeadProvider
}

type Service struct {
	Store
}

func New(store Store) *Service {
	return &Service{
		store,
	}
}
