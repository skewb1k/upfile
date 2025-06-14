package indexFs

import (
	"encoding/base64"
	"fmt"
	"path/filepath"

	"github.com/skewb1k/upfile/pkg/sha256"
)

func (s IndexFsProvider) getEntries() string {
	return filepath.Join(
		s.BaseDir,
		"entries",
	)
}

func (s IndexFsProvider) getPathToEntriesByName(fname string) string {
	return filepath.Join(s.getEntries(), "by-filename", encodePath(fname))
}

func (s IndexFsProvider) getPathToFilenamesByEntry(entry string) string {
	return filepath.Join(s.getEntries(), "by-entry", sha256.FromString(entry).String())
}

func (s IndexFsProvider) getUpstreams() string {
	return filepath.Join(
		s.BaseDir,
		"upstreams",
	)
}

func (s IndexFsProvider) getPathToUpstream(fname string) string {
	return filepath.Join(s.getUpstreams(), encodePath(fname))
}

func encodePath(path string) string {
	return base64.URLEncoding.EncodeToString([]byte(path))
}

func decodePath(encoded string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("failed to decode path: %w", err)
	}

	return string(data), nil
}
