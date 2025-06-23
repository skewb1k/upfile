package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/index"
)

func Diff(
	ctx context.Context,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
	indexProvider IndexProvider,
	path string,
) error {
	path = filepath.Clean(path)
	// TODO: use for built-in pager, now just for checking permissions
	_, err := os.ReadFile(path)
	if err != nil {
		return errors.Unwrap(err)
	}

	fname := filepath.Base(path)

	upstream, err := indexProvider.GetUpstream(ctx, fname)
	if err != nil {
		if errors.Is(err, index.ErrNotFound) {
			return ErrNotTracked
		}

		return fmt.Errorf("get upstream: %w", err)
	}

	tmpFile, err := os.CreateTemp("", "upfile-diff-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(upstream.Content); err != nil {
		return fmt.Errorf("write to temp file: %w", err)
	}
	_ = tmpFile.Close()

	// TODO: do not use git pager or at least have fallback to some built-in one
	gitdiff := exec.Command("git", "diff", "--no-index", tmpFile.Name(), path)
	gitdiff.Stdin = stdin
	gitdiff.Stdout = stdout
	gitdiff.Stderr = stderr

	_ = gitdiff.Run()

	return nil
}
