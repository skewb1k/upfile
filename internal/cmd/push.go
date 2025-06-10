package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/entries"
	"github.com/skewb1k/upfile/internal/upstreams"
	"github.com/spf13/cobra"
)

func pushCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "push <path>",
		Short: "Push file to the upstream",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("failed to get abs path to file: %w", err)
			}

			entryDir, fname := filepath.Dir(path), filepath.Base(path)

			baseDir := getBaseDir()
			upstreamsProvider := upstreams.NewProvider(baseDir)
			entriesProvider := entries.NewProvider(baseDir)

			exists, err := entriesProvider.CheckEntry(cmd.Context(), fname, entryDir)
			if err != nil {
				return fmt.Errorf("check entry: %w", err)
			}

			if !exists {
				return ErrNotTracked
			}

			newContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			upstream, err := upstreamsProvider.GetUpstream(cmd.Context(), fname)
			if err != nil {
				return fmt.Errorf("get upstream: %w", err)
			}

			if upstream.Hash.EqualBytes(newContent) {
				return ErrUpToDate
			}

			if err := upstreamsProvider.SetUpstream(
				cmd.Context(),
				fname,
				upstreams.NewUpstream(string(newContent)),
			); err != nil {
				return fmt.Errorf("set upstream: %w", err)
			}

			return nil
		},
	}
}
