package userfile

import (
	"context"
)

//go:generate go tool mockgen -typed -package userfile -destination ./mock.go . UserFileProvider

type UserFileProvider interface {
	ReadFile(ctx context.Context, path string) (string, error)
}
