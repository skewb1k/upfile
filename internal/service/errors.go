package service

import (
	"errors"
)

var (
	ErrAlreadyTracked  = errors.New("file already tracked")
	ErrCancelled       = errors.New("cancelled")
	ErrNotTracked      = errors.New("not tracked")
	ErrNoEntry         = errors.New("no entry")
	ErrInvalidFilename = errors.New("invalid filename")
	ErrNameUnchanged   = errors.New("names are the same")
	ErrNotDirectory    = errors.New("not a directory")
	ErrDirNotExists    = errors.New("no such directory")
)
