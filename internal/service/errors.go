package service

import (
	"errors"
)

var (
	ErrAlreadyTracked = errors.New("file already tracked")
	ErrUpToDate       = errors.New("up to date")
	ErrNotTracked     = errors.New("file is not tracked")
	ErrNoEntries      = errors.New("no file entries")
	ErrFileNotFound   = errors.New("file not found")
	ErrCancelled      = errors.New("cancelled")
)
