package cmd

import (
	"fmt"
	"path/filepath"

	"upfile/internal/service"
	storeFs "upfile/internal/store/fs"

	"github.com/spf13/cobra"
)

func pull() *cobra.Command {
	var dest string

	cmd := &cobra.Command{
		Use:   "pull <filename>",
		Short: "Pull file from origin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := filepath.Abs(dest)
			if err != nil {
				return fmt.Errorf("get abs path to dest dir: %w", err)
			}

			s := service.New(storeFs.New(getBaseDir()))
			pulled, err := s.Pull(cmd.Context(), args[0], d)
			if err != nil {
				return err
			}

			if !pulled {
				cmd.Println("Already up to date.")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&dest, "dest", "d", ".", "Destination folder")

	return cmd
}
