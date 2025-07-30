package commands

import (
	"fmt"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/service"
	"github.com/spf13/cobra"
)

func Remove() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove <path>",
		Short:   "Remove entry from tracked list",
		Aliases: []string{"rm"},
		Args:    cobra.ExactArgs(1),
		// TODO: provide advanced completions
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("failed to get abs path to file: %w", err)
			}

			if err := service.Remove(cmd.Context(), getIndexFsProvider(), path); err != nil {
				return fmt.Errorf("cannot remove '%s': %w", path, err)
			}

			return nil
		},
	}

	return cmd
}
