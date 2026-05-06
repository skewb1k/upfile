package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func pullCmd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no upstream name given")
	}
	// TODO: support providing multiple upstreams.
	if len(args) > 1 {
		return fmt.Errorf("1 upstream name expected")
	}
	upstream := args[0]
	if !filepath.IsLocal(upstream) {
		return fmt.Errorf("invalid upstream name %s", upstream)
	}

	upfileDir, err := UpfileDir()
	if err != nil {
		return err
	}
	upstreamDir := filepath.Join(upfileDir, "upstream")

	path := filepath.Join(upstreamDir, upstream)
	src, err := os.Open(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("upstream file %s not found", upstream)
		}
		return err
	}
	defer src.Close()

	dst, err := os.Create(upstream)
	if err != nil {
		return err
	}
	defer dst.Close()

	// TODO: perform atomic copy.
	_, err = dst.ReadFrom(src)
	return err
}
