package safejoin

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrEmptyFilename = errors.New("filename is empty")
	ErrContainsSlash = errors.New("filename must not contain slashes")
	ErrIsAbsolute    = errors.New("filename must not be absolute")
	ErrHasVolumeName = errors.New("filename must not contain volume name")
	ErrEscapesBase   = errors.New("path escapes base directory")
)

// SafeJoinFilename safely joins a trusted base directory with an untrusted filename.
// It ensures the result is a path to a file inside baseDir, without subdirectories.
// It rejects empty, absolute, or slash-containing names (including Windows-style).
func SafeJoinFilename(baseDir, fname string) (string, error) {
	if fname == "" {
		return "", ErrEmptyFilename
	}

	normalized := filepath.ToSlash(fname)

	if strings.Contains(normalized, "/") {
		return "", ErrContainsSlash
	}

	if filepath.IsAbs(normalized) {
		return "", ErrIsAbsolute
	}

	if vol := filepath.VolumeName(fname); vol != "" {
		return "", ErrHasVolumeName
	}

	fullPath := filepath.Join(baseDir, normalized)

	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return "", fmt.Errorf("resolve base directory: %w", err)
	}

	absFull, err := filepath.Abs(fullPath)
	if err != nil {
		return "", fmt.Errorf("resolve full path: %w", err)
	}

	if !strings.HasPrefix(absFull, absBase+string(os.PathSeparator)) {
		return "", ErrEscapesBase
	}

	return fullPath, nil
}
