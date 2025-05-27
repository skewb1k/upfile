package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func diff() *cobra.Command {
	return &cobra.Command{
		Use:   "diff <path>",
		Short: "Diff with origin",
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
			_ = baseDir

			// c := commitsFs.NewStore(baseDir)
			// e := entriesFs.NewStore(baseDir)

			// s := service.New(c, e)

			// res, err := s.Diff(cmd.Context(), path)
			// if err != nil {
			// 	return fmt.Errorf("failed to link: %w", err)
			// }

			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Linked: %s\n", path); err != nil {
				return fmt.Errorf("failed to write: %w", err)
			}

			return nil
		},
	}
}
