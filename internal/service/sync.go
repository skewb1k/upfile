package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/index"
)

func Sync(
	ctx context.Context,
	stdin io.Reader,
	stdout io.Writer,
	indexProvider IndexProvider,
	yes bool,
	fname string,
) error {
	entries, err := indexProvider.GetEntriesByFilename(ctx, fname)
	if err != nil {
		return fmt.Errorf("get entries by filename: %w", err)
	}

	if len(entries) == 0 {
		mustFmt(fmt.Fprintln(stdout, "No file entries"))
		return nil
	}

	upstream, err := indexProvider.GetUpstream(ctx, fname)
	if err != nil {
		if errors.Is(err, index.ErrNotFound) {
			return ErrNotTracked
		}

		return fmt.Errorf("get upstream: %w", err)
	}

	toUpdate := make([]string, 0)

	for _, entry := range entries {
		path := filepath.Join(entry, fname)

		existing, err := os.ReadFile(path)
		if err == nil && upstream.Hash.EqualBytes(existing) {
			// Up-to-date, skip
			continue
		}

		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}

		toUpdate = append(toUpdate, path)
	}

	if len(toUpdate) == 0 {
		mustFmt(fmt.Fprintln(stdout, "All entries up-to-date"))
		return nil
	}

	if !yes {
		mustFmt(fmt.Fprintln(stdout, "The following entries will be updated:"))

		for _, e := range toUpdate {
			mustFmt(fmt.Fprintln(stdout, " -", e))
		}

		ok, err := promptDefaultYes(stdin, stdout)
		if err != nil {
			return err
		}

		if !ok {
			return ErrCancelled
		}
	}

	for _, fullPath := range toUpdate {
		if err := MkdirAllWriteFile(fullPath, upstream.Content); err != nil {
			return err
		}
	}

	return nil
}
