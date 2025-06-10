package upstreams

import (
	"encoding/base64"
	"fmt"
	"path/filepath"
)

func (p Provider) getUpstreams() string {
	return filepath.Join(
		p.BaseDir,
		"upstreams",
	)
}

func (p Provider) getPathToUpstream(fname string) string {
	return filepath.Join(p.getUpstreams(), encodePath(fname))
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
