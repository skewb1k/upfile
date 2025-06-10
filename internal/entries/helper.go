package entries

import (
	"encoding/base64"
	"path/filepath"

	"github.com/skewb1k/upfile/pkg/sha256"
)

func (p Provider) getEntries() string {
	return filepath.Join(
		p.BaseDir,
		"entries",
	)
}

func (p Provider) getPathToEntriesByName(fname string) string {
	return filepath.Join(p.getEntries(), "by-filename", encodePath(fname))
}

func (p Provider) getPathToFilenamesByEntry(entry string) string {
	return filepath.Join(p.getEntries(), "by-entry", sha256.FromString(entry).String())
}

func encodePath(path string) string {
	return base64.URLEncoding.EncodeToString([]byte(path))
}
