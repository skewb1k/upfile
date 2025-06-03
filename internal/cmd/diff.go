package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"upfile/internal/service"

	"github.com/spf13/cobra"
)

func diff() *cobra.Command {
	return &cobra.Command{
		Use:   "diff <path>",
		Short: "Diff with origin",
		Args:  cobra.ExactArgs(1),
		RunE: wrap(func(cmd *cobra.Command, s *service.Service, args []string) error {
			path, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("get abs path to file: %w", err)
			}

			fname := filepath.Base(path)

			upstreamContent, err := s.Diff(cmd.Context(), fname)
			if err != nil {
				return err //nolint: wrapcheck
			}

			return gitDiff(cmd.OutOrStdout(), cmd.ErrOrStderr(), path, upstreamContent)
		}),
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
