package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"upfile/internal/commits"
	"upfile/internal/entries"
	"upfile/internal/head"
)

type Service struct {
	commits commits.CommitStore
	head    head.HeadStore
	entries entries.EntryStore
}

func New(
	commitStore commits.CommitStore,
	head head.HeadStore,
	entryStore entries.EntryStore,
) *Service {
	return &Service{
		commits: commitStore,
		head:    head,
		entries: entryStore,
	}
}

func computeHash(data []byte) string {
	checksum := sha256.Sum256(data)
	return hex.EncodeToString(checksum[:])
}

var ErrAlreadyLinked = errors.New("file already linked")
