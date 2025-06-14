package service

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/skewb1k/upfile/internal/index"
	"github.com/skewb1k/upfile/pkg/validfname"
)

func Show(
	ctx context.Context,
	stdout io.Writer,
	indexProvider IndexProvider,
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

	mustFmt(fmt.Fprint(stdout, string(upstream.Content)))

	return nil
}
