package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"upfile/internal/service"
	storeFs "upfile/internal/store/fs"

	"github.com/spf13/cobra"
)

func FileExistsAndReadable(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	_ = f.Close()
	return nil
}

func diff() *cobra.Command {
	return &cobra.Command{
		Use:   "diff <path>",
		Short: "Diff with origin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("get abs path to file: %w", err)
			}

			if err := FileExistsAndReadable(path); err != nil {
				return err
			}

			s := service.New(storeFs.New(getBaseDir()))

			upstreamContent, err := s.GetUpstream(cmd.Context(), filepath.Base(args[0]))
			if err != nil {
				return err
			}

			return gitDiff(cmd.OutOrStdout(), cmd.ErrOrStderr(), path, upstreamContent)
		},
	}
}

func gitDiff(stdout, stderr io.Writer, filePath string, content string) error {
	tmpFile, err := os.CreateTemp("", Name+"-diff-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		return fmt.Errorf("write to temp file: %w", err)
	}
	_ = tmpFile.Close()

	cmd := exec.Command("git", "diff", "--no-index", tmpFile.Name(), filePath)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	_ = cmd.Run()

	return nil
}
