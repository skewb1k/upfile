package cmd

import (
	"fmt"
	"path/filepath"

	"upfile/internal/service"
	storeFs "upfile/internal/store/fs"

	"github.com/spf13/cobra"
)

func add() *cobra.Command {
	return &cobra.Command{
		Use:   "add <path>",
		Short: "Add a file to be tracked",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := filepath.Abs(filepath.Clean(args[0]))
			if err != nil {
				return fmt.Errorf("get abs path to file: %w", err)
			}

			s := service.New(storeFs.New(getBaseDir()))

			if err := s.Add(cmd.Context(), path); err != nil {
				return err
			}

			cmd.Printf("Added: %s\n", path)

			return nil
		},
	}
}
