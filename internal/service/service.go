package service

import "context"

type IndexProvider interface {
	CheckEntry(ctx context.Context, fname string, entry string) (bool, error)
	CheckUpstream(ctx context.Context, fname string) (bool, error)
	CreateEntry(ctx context.Context, fname string, entry string) error
	DeleteEntry(ctx context.Context, fname string, entry string) error
	DeleteUpstream(ctx context.Context, fname string) error
	GetEntriesByFilename(ctx context.Context, fname string) ([]string, error)
	GetFilenames(ctx context.Context) ([]string, error)
	GetFilenamesByEntry(ctx context.Context, entry string) ([]string, error)
	GetUpstream(ctx context.Context, fname string) (Upstream, error)
	SetUpstream(ctx context.Context, fname string, upstream *Upstream) error
}
