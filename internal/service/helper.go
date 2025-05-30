package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

func computeHash(data []byte) string {
	checksum := sha256.Sum256(data)
	return hex.EncodeToString(checksum[:])
}

var (
	ErrAlreadyLinked = errors.New("file already linked")
	ErrNotTracked    = errors.New("file is not tracked")
)
