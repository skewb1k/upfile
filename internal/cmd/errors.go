package cmd

import (
	"errors"
)

var (
	ErrAlreadyTracked = errors.New("file already tracked")
	ErrNotTracked     = errors.New("file is not tracked")
	ErrFileNotFound   = errors.New("file not found")
	ErrNoEntries      = errors.New("no file entries")
)
