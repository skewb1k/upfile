package entries

import "errors"

var (
	ErrExists          = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
	ErrInvalidFilename = errors.New("invalid filename")
)
