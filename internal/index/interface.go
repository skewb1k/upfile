package index

import (
	"context"
)

//go:generate go tool mockgen -typed -package index -destination ./mock.go . IndexProvider

type IndexProvider interface {
	GetFiles(ctx context.Context) ([]string, error)
	GetFilesByEntryDir(ctx context.Context, entryDir string) ([]string, error)

	CreateEntry(ctx context.Context, fname string, entryDir string) error
	CheckEntry(ctx context.Context, fname string, entryDir string) (bool, error)
	GetEntriesByFname(ctx context.Context, fname string) ([]string, error)
	DeleteEntry(ctx context.Context, fname string, entryDir string) error

	SetUpstream(ctx context.Context, fname string, value string) error
	GetUpstream(ctx context.Context, fname string) (string, error)
	CheckUpstream(ctx context.Context, fname string) (bool, error)
}
