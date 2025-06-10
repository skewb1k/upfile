package upstreams

import (
	"github.com/skewb1k/upfile/pkg/sha256"
)

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
