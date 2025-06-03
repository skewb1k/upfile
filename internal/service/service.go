package service

import (
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
