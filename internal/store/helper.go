package store

import (
	"encoding/base64"
	"path/filepath"

	"github.com/skewb1k/upfile/pkg/sha256"
)

func (p Store) getEntries() string {
	return filepath.Join(
		p.BaseDir,
		"index",
	)
}

func (p Store) getPathToEntriesByName(fname string) string {
	return filepath.Join(p.getEntries(), "by-filename", encodePath(fname))
}

func (p Store) getPathToFilenamesByEntry(entry string) string {
	return filepath.Join(p.getEntries(), "by-entry", sha256.FromString(entry).Hex())
}

func encodePath(path string) string {
	return base64.URLEncoding.EncodeToString([]byte(path))
}

func (p Store) getHeadsPath() string {
	return filepath.Join(
		p.BaseDir,
		"heads",
	)
}

func (p Store) getPathToCommits() string {
	return filepath.Join(
		p.BaseDir,
		"commits",
	)
}

func (p Store) getPathToBlobs() string {
	return filepath.Join(
		p.BaseDir,
		"blobs",
	)
}
