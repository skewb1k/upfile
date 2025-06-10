package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/upstreams"
	"github.com/spf13/cobra"
)

func diffCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diff <path>",
		Short: "Show the difference between a file and its upstream version",
		Long: doc(`
This command compares the contents of the given file (by absolute or relative path)
to its upstream.

You must pass a path to a file that is already being tracked. The file is matched by its filename
against the known tracked entries.`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("failed to get abs path to file: %w", err)
			}

			fname := filepath.Base(path)

			baseDir := getBaseDir()
			upstreamsProvider := upstreams.NewProvider(baseDir)

			upstream, err := upstreamsProvider.GetUpstream(cmd.Context(), fname)
			if err != nil {
				return fmt.Errorf("get upstream: %w", err)
			}

			tmpFile, err := os.CreateTemp("", name+"-diff-*.tmp")
			if err != nil {
				return fmt.Errorf("create temp file: %w", err)
			}
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.WriteString(upstream.Content); err != nil {
				return fmt.Errorf("write to temp file: %w", err)
			}
			_ = tmpFile.Close()

			// TODO: do not use git pager or at least have fallback to some built-in one
			gitdiff := exec.Command("git", "diff", "--no-index", tmpFile.Name(), path)
			gitdiff.Stdout = cmd.OutOrStdout()
			gitdiff.Stderr = cmd.ErrOrStderr()

			_ = gitdiff.Run()

			return nil
		},
	}

	return cmd
}
