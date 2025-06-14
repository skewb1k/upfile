package cmd

import (
	"errors"
)

var (
	ErrAlreadyTracked  = errors.New("file already tracked")
	ErrNotTracked      = errors.New("file is not tracked")
	ErrInvalidFilename = errors.New("invalid filename")
	ErrFileNotFound    = errors.New("file not found")
	ErrNoEntries       = errors.New("no file entries")
	ErrNameUnchanged   = errors.New("old and new names are the same")
)
