package commands

import (
	"fmt"
	"os"
	"path/filepath"

	indexFs "github.com/skewb1k/upfile/internal/index/fs"
)

func getBaseDir() string {
	if dir := os.Getenv("UPFILE_DIR"); dir != "" {
		return dir
	}

	if xdgData := os.Getenv("XDG_DATA_HOME"); xdgData != "" {
		return filepath.Join(xdgData, "upfile")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("failed to get current user's home dir: %s", err))
	}

	return filepath.Join(home, ".local", "share", "upfile")
}

func getIndexFsProvider() *indexFs.IndexFsProvider {
	return indexFs.NewProvider(getBaseDir())
}

func doc(s string) string {
	return s[1:]
}
