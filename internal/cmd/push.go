package cmd

import (
	"context"
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

func pushCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push <path>",
		Short: "Push file to the upstream",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("failed to get abs path to file: %w", err)
			}

			return Push(cmd.Context(), service.New(store.New(getBaseDir())), path)
		},
	}

	return cmd
}

func Push(
	ctx context.Context,
	s *service.Service,
	path string,
) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	fname, entry := filepath.Base(path), filepath.Dir(path)

	entryExists, err := s.CheckEntry(ctx, fname, entry)
	if err != nil {
		return fmt.Errorf("create entry: %w", err)
	}

	if !entryExists {
		return ErrNotTracked
	}

	head, err := s.GetHead(ctx, fname)
	if err != nil {
		return fmt.Errorf("get head: %w", err)
	}

	headCommit, err := s.GetCommit(ctx, head)
	if err != nil {
		return fmt.Errorf("get commit: %w", err)
	}

	contentHash := sha256.FromBytes(content)

	if headCommit.ContentHash == contentHash {
		return ErrUpToDate
	}

	commit := &commit.Commit{
		Filename:    fname,
		ContentHash: contentHash,
		Parent:      &head,
		Timestamp:   time.Now(),
	}

	commitHash := commit.HashCommit()

	if err := s.SaveCommit(ctx, commitHash, commit); err != nil {
		return fmt.Errorf("save commit: %w", err)
	}

	return nil
}
