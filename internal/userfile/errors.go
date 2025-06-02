package userfile

import "errors"

var (
	ErrExists   = errors.New("already exists")
	ErrNotFound = errors.New("not found")
)
