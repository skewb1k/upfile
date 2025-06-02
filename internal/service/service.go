package service

import (
	"upfile/internal/index"
	"upfile/internal/userfile"
)

type Service struct {
	indexStore    index.Store
	userfileStore userfile.Store
}

func New(
	indexStore index.Store,
	userfileStore userfile.Store,
) *Service {
	return &Service{
		indexStore:    indexStore,
		userfileStore: userfileStore,
	}
}
