package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/index"
	"github.com/skewb1k/upfile/pkg/validfname"
)

func Pull(
	ctx context.Context,
	stdin io.Reader,
	stdout io.Writer,
	indexProvider IndexProvider,
	yes bool,
	dest string,
	fname string,
) error {
	if !validfname.ValidateFilename(fname) {
		return ErrInvalidFilename
	}

	upstream, err := indexProvider.GetUpstream(ctx, fname)
	if err != nil {
		if errors.Is(err, index.ErrNotFound) {
			return ErrNotTracked
		}

		return fmt.Errorf("get upstream: %w", err)
	}

	path := filepath.Join(dest, fname)

	existing, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	} else {
		if upstream.Hash.EqualBytes(existing) {
			mustFmt(fmt.Fprintln(stdout, "File up-to-date"))
			return nil
		}

		if !yes {
			mustFmt(fmt.Fprintf(stdout, "File '%s' will be updated\n", path))

			ok, err := promptDefaultYes(stdin, stdout)
			if err != nil {
				return err
			}

			if !ok {
				os.Exit(1)
			}
		}
	}

	if err := indexProvider.CreateEntry(ctx, fname, dest); err != nil {
		if !errors.Is(err, index.ErrExists) {
			return fmt.Errorf("create entry: %w", err)
		}
	}

	if err := MkdirAllWriteFile(path, upstream.Content); err != nil {
		return err
	}

	return nil
}
