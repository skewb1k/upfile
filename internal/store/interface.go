package store

import (
	"context"
)

type Commit struct {
	Hash    string
	Content string
	Parent  string
}

type Store interface {
	GetFiles(ctx context.Context) ([]string, error)

	CreateCommit(ctx context.Context, fname string, commit *Commit) error
	GetCommitByHash(ctx context.Context, fname string, hash string) (Commit, error)

	CreateEntry(ctx context.Context, fname string, entry string) error
	CheckEntry(ctx context.Context, fname string, entry string) (bool, error)
	GetEntries(ctx context.Context, fname string) ([]string, error)
	DeleteEntry(ctx context.Context, fname string, entry string) error

	SetHead(ctx context.Context, fname string, value string) error
	SetHeadIfNotExists(ctx context.Context, fname string, value string) error
	GetHead(ctx context.Context, fname string) (string, error)
}
