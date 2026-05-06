package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func pushCmd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no file given")
	}
	// TODO: support providing multiple paths.
	if len(args) > 1 {
		return fmt.Errorf("1 file expected")
	}
	srcPath := args[0]

	upfileDir, err := UpfileDir()
	if err != nil {
		return err
	}
	upstreamDir := filepath.Join(upfileDir, "upstream")
	if err := os.MkdirAll(upstreamDir, 0o700); err != nil {
		return err
	}

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	upstream := filepath.Base(srcPath)
	upstreamPath := filepath.Join(upstreamDir, upstream)
	upstreamFile, err := os.OpenFile(upstreamPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer upstreamFile.Close()

	// TODO: perform atomic copy.
	_, err = upstreamFile.ReadFrom(srcFile)
	return err
}
