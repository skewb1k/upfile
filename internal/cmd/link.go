package cmd

import (
	"fmt"
	"path/filepath"

	"upfile/internal/service"
	storeFs "upfile/internal/store/fs"

	"github.com/spf13/cobra"
)

func link() *cobra.Command {
	return &cobra.Command{
		Use:   "link <path>",
		Short: "Link a file to be tracked",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := filepath.Abs(filepath.Clean(args[0]))
			if err != nil {
				return fmt.Errorf("failed to get abs path to file: %w", err)
			}

			baseDir := getBaseDir()

			if err := service.Link(cmd.Context(),
				storeFs.New(baseDir),
				path,
			); err != nil {
				return fmt.Errorf("failed to link: %w", err)
			}

			cmd.Printf("Linked: %s\n", path)

			return nil
		},
	}
}
