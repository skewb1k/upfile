package index

import (
	"context"
)

//go:generate go tool mockgen -typed -package index -destination ./mock.go . IndexProvider

type IndexProvider interface {
	GetFilenames(ctx context.Context) ([]string, error)
	GetFilenamesByEntry(ctx context.Context, entry string) ([]string, error)

	CreateEntry(ctx context.Context, fname string, entry string) error
	CheckEntry(ctx context.Context, fname string, entry string) (bool, error)
	DeleteEntry(ctx context.Context, fname string, entry string) error
	GetEntriesByFilename(ctx context.Context, fname string) ([]string, error)

	SetUpstream(ctx context.Context, fname string, value string) error
	GetUpstream(ctx context.Context, fname string) (string, error)
	CheckUpstream(ctx context.Context, fname string) (bool, error)
	DeleteUpstream(ctx context.Context, fname string) error
}
