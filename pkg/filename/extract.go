package filename

import (
	"path/filepath"
)

func Extract(path string) (string, bool) {
	fname := filepath.Base(filepath.Clean(path))

	switch fname {
	case ".", "/", "..":
		return "", false
	}

	return fname, true
}
