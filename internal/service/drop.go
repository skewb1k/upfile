package service

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/skewb1k/upfile/internal/index"
)

func Drop(
	ctx context.Context,
	stdin io.Reader,
	stdout io.Writer,
	indexProvider IndexProvider,
	yes bool,
	fname string,
) error {
	upstreamExists, err := indexProvider.CheckUpstream(ctx, fname)
	if err != nil {
		return fmt.Errorf("check upstream: %w", err)
	}

	if !upstreamExists {
		return ErrNotTracked
	}

	entries, err := indexProvider.GetEntriesByFilename(ctx, fname)
	if err != nil {
		return fmt.Errorf("get entries by filename: %w", err)
	}

	if !yes {
		if len(entries) == 0 {
			mustFmt(fmt.Fprintf(stdout, "'%s' upstream will be deleted\n", fname))
		} else {
			mustFmt(fmt.Fprintf(stdout, "'%s' upstream will be deleted and following entries will be untracked:\n", fname))
			for _, e := range entries {
				mustFmt(fmt.Fprintln(stdout, " -", e))
			}
		}

		ok, err := promptDefaultNo(stdin, stdout)
		if err != nil {
			return err
		}

		if !ok {
			return ErrCancelled
		}
	}

	if err := indexProvider.DeleteUpstream(ctx, fname); err != nil {
		if errors.Is(err, index.ErrNotFound) {
			return ErrNotTracked
		}

		return fmt.Errorf("delete upstream: %w", err)
	}

	for _, entry := range entries {
		if err := indexProvider.DeleteEntry(ctx, fname, entry); err != nil {
			return fmt.Errorf("delete entry: %w", err)
		}
	}

	return nil
}
