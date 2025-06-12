package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/skewb1k/upfile/internal/commit"
	"github.com/skewb1k/upfile/internal/service"
	"github.com/skewb1k/upfile/internal/store"
	"github.com/skewb1k/upfile/pkg/sha256"
	"github.com/spf13/cobra"
)

func addCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <path>",
		Short: "Add a file to be tracked",
		Args:  cobra.ExactArgs(1),
		RunE: wrap(func(cmd *cobra.Command, s *service.Service, args []string) error {
			path, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("failed to get abs path to file: %w", err)
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			fname, entry := filepath.Base(path), filepath.Dir(path)

			if err := s.CreateEntry(cmd.Context(), fname, entry); err != nil {
				if errors.Is(err, store.ErrExists) {
					return ErrAlreadyTracked
				}

				return fmt.Errorf("create entry: %w", err)
			}

			headExists, err := s.CheckHead(cmd.Context(), fname)
			if err != nil {
				return fmt.Errorf("check head: %w", err)
			}

			if !headExists {
				// file was not tracked, create first commit

				contentHash := sha256.FromBytes(content)
				if err := s.SaveBlob(cmd.Context(), contentHash, content); err != nil {
					// TODO: handle conflict
					return fmt.Errorf("save blob: %w", err)
				}

				c := &commit.Commit{
					Filename:    fname,
					ContentHash: contentHash,
					Parent:      nil,
					Timestamp:   time.Now(),
				}
				commitHash := c.HashCommit()

				if err := s.SaveCommit(cmd.Context(), commitHash, c); err != nil {
					return fmt.Errorf("save commit: %w", err)
				}

				if err := s.SetHead(cmd.Context(), fname, commitHash); err != nil {
					return fmt.Errorf("set head: %w", err)
				}
			}

			return nil
		}),
	}

	return cmd
}
