package service

import (
	"crypto/sha256"

	"upfile/internal/index"
	"upfile/internal/userfile"
)

type Service struct {
	indexProvider    index.IndexProvider
	userfileProvider userfile.UserFileProvider
}

func New(
	indexProvider index.IndexProvider,
	userfileProvider userfile.UserFileProvider,
) *Service {
	return &Service{
		indexProvider:    indexProvider,
		userfileProvider: userfileProvider,
	}
}

func hash(s string) [32]byte {
	return sha256.Sum256([]byte(s))
}
