package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/service"
	"github.com/spf13/cobra"
)

func syncCmd() *cobra.Command {
	// TODO:
	// var patch bool

	var yes bool

	cmd := &cobra.Command{
		Use:   "sync <filename>",
		Args:  cobra.ExactArgs(1),
		Short: "Sync all entries of file with upstream",
		RunE: wrap(func(cmd *cobra.Command, s *service.Service, args []string) error {
			fname := args[0]

			entriesList, err := s.GetEntriesByFilename(cmd.Context(), fname)
			if err != nil {
				return fmt.Errorf("get entries by filename: %w", err)
			}

			head, err := s.GetHead(cmd.Context(), fname)
			if err != nil {
				return fmt.Errorf("get head: %w", err)
			}

			headCommit, err := s.GetCommit(cmd.Context(), head)
			if err != nil {
				return fmt.Errorf("get commit: %w", err)
			}

			content, err := s.GetBlob(cmd.Context(), headCommit.ContentHash)
			if err != nil {
				return fmt.Errorf("get blob: %w", err)
			}

			toUpdate := make([]string, 0)

			for _, entryDir := range entriesList {
				path := filepath.Join(entryDir, fname)

				existing, err := os.ReadFile(path)
				if err != nil {
					return err
				}

				if !headCommit.ContentHash.EqualBytes(existing) {
					toUpdate = append(toUpdate, filepath.Join(entryDir, fname))
				}
			}

			if len(toUpdate) == 0 {
				return ErrUpToDate
			}

			if !yes && !ask(cmd.InOrStdin(), toUpdate, true, "The following entries will be updated:") {
				os.Exit(1)
				return nil
			}

			for _, fullPath := range toUpdate {
				if err := WriteFile(fullPath, content); err != nil {
					return fmt.Errorf("write file: %w", err)
				}
			}

			return nil
		}),
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Automatic 'yes' to prompts")
	cmd.ValidArgsFunction = completeFname

	// cmd.Flags().BoolVarP(&patch, "patch", "p", false, "Interactively apply changes per entry")

	return cmd
}
