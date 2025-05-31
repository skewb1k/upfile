package cmd

import (
	"fmt"
	"path/filepath"

	"upfile/internal/service"
	storeFs "upfile/internal/store/fs"

	"github.com/spf13/cobra"
)

func remove() *cobra.Command {
	return &cobra.Command{
		Use:     "remove <path>",
		Short:   "Remove entry from tracked list",
		Aliases: []string{"rm"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := filepath.Abs(filepath.Clean(args[0]))
			if err != nil {
				return fmt.Errorf("get abs path to file: %w", err)
			}

			s := service.New(storeFs.New(getBaseDir()))

			if err := s.Remove(cmd.Context(), path); err != nil {
				return err
			}

			cmd.Printf("Removed: %s\n", path)

			return nil
		},
	}
}
