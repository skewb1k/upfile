package entries

import (
	"context"
	"errors"
)

type EntryStore interface {
	CreateEntry(ctx context.Context, fname string, entry string) error
	CheckEntry(ctx context.Context, fname string, entry string) (bool, error)
	GetEntries(ctx context.Context, fname string) ([]string, error)
}

var (
	ErrExists   = errors.New("already exists")
	ErrNotFound = errors.New("not found")
)
