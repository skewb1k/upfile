package userfile

import (
	"context"
)

//go:generate go tool mockgen -typed -package userfile -destination ./store_mock.go . Store

type Store interface {
	ReadFile(ctx context.Context, path string) (string, error)
}
