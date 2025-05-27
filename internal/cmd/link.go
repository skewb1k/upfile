package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	commitsFs "upfile/internal/commits/fs"
	entriesFs "upfile/internal/entries/fs"
	headFs "upfile/internal/head/fs"
	"upfile/internal/service"

	"github.com/spf13/cobra"
)

const globalDir = ".local/upfile"

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

			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get current user's home dir: %w", err)
			}

			baseDir := filepath.Join(home, globalDir)

			c := commitsFs.NewStore(baseDir)
			h := headFs.NewStore(baseDir)
			e := entriesFs.NewStore(baseDir)

			s := service.New(c, h, e)

			if err := s.Link(cmd.Context(), path); err != nil {
				if errors.Is(err, service.ErrAlreadyLinked) {
					cmd.Println("Already linked")
					return nil
				}

				return fmt.Errorf("failed to link: %w", err)
			}

			cmd.Printf("Linked: %s\n", path)

			return nil
		},
	}
}
