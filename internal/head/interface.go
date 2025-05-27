package head

import (
	"context"
	"errors"
)

type HeadStore interface {
	SetHead(ctx context.Context, fname string, value string) error
	GetHead(ctx context.Context, fname string) (string, error)
}

var ErrNotFound = errors.New("not found")
