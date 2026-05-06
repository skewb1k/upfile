package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func listCmd(args []string) error {
	upfileDir, err := UpfileDir()
	if err != nil {
		return err
	}
	upstreamDir := filepath.Join(upfileDir, "upstream")
	upstreams, err := os.ReadDir(upstreamDir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}
	for _, upstream := range upstreams {
		if upstream.IsDir() {
			continue
		}
		fmt.Println(upstream.Name())
	}
	return nil
}
