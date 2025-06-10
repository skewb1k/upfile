package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/entries"
	"github.com/skewb1k/upfile/internal/upstreams"
	"github.com/spf13/cobra"
)

func addCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <path>",
		Short: "Add a file to be tracked",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("failed to get abs path to file: %w", err)
			}

			baseDir := getBaseDir()
			upstreamsProvider := upstreams.NewProvider(baseDir)
			entriesProvider := entries.NewProvider(baseDir)

			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			fname, entryDir := filepath.Base(path), filepath.Dir(path)

			if err := entriesProvider.CreateEntry(cmd.Context(), fname, entryDir); err != nil {
				return fmt.Errorf("create entry: %w", err)
			}

			upstreamExists, err := upstreamsProvider.CheckUpstream(cmd.Context(), fname)
			if err != nil {
				return fmt.Errorf("check upstream: %w", err)
			}

			if !upstreamExists {
				if err := upstreamsProvider.SetUpstream(cmd.Context(), fname, upstreams.NewUpstream(string(content))); err != nil {
					return fmt.Errorf("set upstream: %w", err)
				}
			}

			return nil
		},
	}

	return cmd
}
