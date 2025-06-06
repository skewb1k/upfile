package service

import (
	"errors"
)

var (
	ErrAlreadyTracked = errors.New("file already tracked")
	ErrNotTracked     = errors.New("file is not tracked")
	ErrFileNotFound   = errors.New("file not found")
	ErrUpToDate       = errors.New("up to date")
	ErrNoEntries      = errors.New("no file entries")
	ErrCancelled      = errors.New("cancelled")
)
