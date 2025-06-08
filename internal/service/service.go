package service

import (
	"github.com/skewb1k/upfile/internal/index"
	"github.com/skewb1k/upfile/internal/userfile"
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
