package index

import (
	"context"

	"github.com/skewb1k/upfile/pkg/sha256"
)

//go:generate go tool mockgen -typed -package index -destination ./mock.go . IndexProvider

type IndexProvider interface {
	GetFilenames(ctx context.Context) ([]string, error)
	GetFilenamesByEntry(ctx context.Context, entry string) ([]string, error)

	CreateEntry(ctx context.Context, fname string, entry string) error
	CheckEntry(ctx context.Context, fname string, entry string) (bool, error)
	DeleteEntry(ctx context.Context, fname string, entry string) error
	GetEntriesByFilename(ctx context.Context, fname string) ([]string, error)

	SetUpstream(ctx context.Context, fname string, upstream *Upstream) error
	GetUpstream(ctx context.Context, fname string) (Upstream, error)
	CheckUpstream(ctx context.Context, fname string) (bool, error)
	DeleteUpstream(ctx context.Context, fname string) error
}

type Upstream struct {
	Hash    sha256.SHA256
	Content string
}

func NewUpstream(content string) *Upstream {
	return &Upstream{
		Hash:    sha256.FromString(content),
		Content: content,
	}
}
