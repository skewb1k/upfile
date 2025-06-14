package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/skewb1k/upfile/pkg/validfname"
)

func Rename(
	ctx context.Context,
	indexProvider IndexProvider,
	oldName string,
	newName string,
) error {
	if !validfname.ValidateFilename(newName) {
		return ErrInvalidFilename
	}

	if oldName == newName {
		return ErrNameUnchanged
	}

	oldNameUpstreamExists, err := indexProvider.CheckUpstream(ctx, oldName)
	if err != nil {
		return fmt.Errorf("check upstream: %w", err)
	}

	if !oldNameUpstreamExists {
		return ErrNotTracked
	}

	// Check if new name is not already tracked
	newNameUpstreamExists, err := indexProvider.CheckUpstream(ctx, newName)
	if err != nil {
		return fmt.Errorf("check upstream: %w", err)
	}

	if newNameUpstreamExists {
		return ErrAlreadyTracked
	}

	entries, err := indexProvider.GetEntriesByFilename(ctx, oldName)
	if err != nil {
		return fmt.Errorf("get entries by filename: %w", err)
	}

	// Check for conflicts
	for _, dir := range entries {
		newPath := filepath.Join(dir, newName)
		if _, err := os.Stat(newPath); err == nil {
			return fmt.Errorf("target path %s: %w", newPath, ErrTargetPathAlreadyExists)
		} else if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("stat target path %s: %w", newPath, err)
		}
	}

	// Perform renames and index updates
	for _, dir := range entries {
		oldPath := filepath.Join(dir, oldName)
		newPath := filepath.Join(dir, newName)

		if err := os.Rename(oldPath, newPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("rename file in dir %s: %w", dir, err)
		}

		if err := indexProvider.CreateEntry(ctx, newName, dir); err != nil {
			return fmt.Errorf("add new entry: %w", err)
		}
		if err := indexProvider.DeleteEntry(ctx, oldName, dir); err != nil {
			return fmt.Errorf("delete old entry: %w", err)
		}
	}

	// Copy upstream
	upstream, err := indexProvider.GetUpstream(ctx, oldName)
	if err != nil {
		return fmt.Errorf("get upstream: %w", err)
	}

	if err := indexProvider.SetUpstream(ctx, newName, &upstream); err != nil {
		return fmt.Errorf("set upstream: %w", err)
	}
	if err := indexProvider.DeleteUpstream(ctx, oldName); err != nil {
		return fmt.Errorf("delete old upstream: %w", err)
	}

	return nil
}
