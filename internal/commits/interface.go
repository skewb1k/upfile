package commits

import (
	"context"
	"errors"
)

type Commit struct {
	Hash    string
	Content string
	Parent  string
}

type CommitStore interface {
	CreateCommit(ctx context.Context, fname string, commit *Commit) error
	GetCommitByHash(ctx context.Context, fname string, hash string) (Commit, error)
}

var ErrNotFound = errors.New("not found")
