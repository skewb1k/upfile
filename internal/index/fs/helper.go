package indexFs

import (
	"encoding/base64"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/index"
	"github.com/skewb1k/upfile/pkg/safejoin"
)

func (p Provider) getEntries() string {
	return filepath.Join(
		p.BaseDir,
		"entries",
	)
}

func (p Provider) getPathToEntriesByName(fname string) (string, error) {
	base := filepath.Join(p.getEntries(), "by-filename")

	res, err := safejoin.SafeJoinFilename(base, fname)
	if err != nil {
		return "", index.ErrInvalidFilename
	}

	return res, nil
}

func (p Provider) getPathToFilenamesByEntry(entry string) string {
	return filepath.Join(p.getEntries(), "by-entry", entry)
}

func encodePath(path string) string {
	return base64.URLEncoding.EncodeToString([]byte(path))
}

func (p Provider) getUpstreams() string {
	return filepath.Join(
		p.BaseDir,
		"upstreams",
	)
}

func (p Provider) getPathToUpstream(fname string) (string, error) {
	dir := p.getUpstreams()

	path, err := safejoin.SafeJoinFilename(dir, fname)
	if err != nil {
		return "", index.ErrInvalidFilename
	}

	return path, nil
}
