package service

import "github.com/skewb1k/upfile/pkg/sha256"

type Upstream struct {
	Hash    sha256.SHA256
	Content []byte
}

func New(content []byte) *Upstream {
	return &Upstream{
		Hash:    sha256.FromBytes(content),
		Content: content,
	}
}
