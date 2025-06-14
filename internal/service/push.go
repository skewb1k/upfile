package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Push(
	ctx context.Context,
	stdout io.Writer,
	indexProvider IndexProvider,
	path string,
) error {
	newContent, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	entry, fname := filepath.Dir(path), filepath.Base(path)

	exists, err := indexProvider.CheckEntry(ctx, fname, entry)
	if err != nil {
		return fmt.Errorf("check entry: %w", err)
	}

	if !exists {
		return ErrNoEntry
	}

	upstrem, err := indexProvider.GetUpstream(ctx, fname)
	if err != nil {
		return fmt.Errorf("get upstream: %w", err)
	}

	if upstrem.Hash.EqualBytes(newContent) {
		mustFmt(fmt.Fprintln(stdout, "File up-to-date"))
		return nil
	}

	if err := indexProvider.SetUpstream(
		ctx,
		fname,
		New(newContent),
	); err != nil {
		return fmt.Errorf("set upstream: %w", err)
	}

	return nil
}
