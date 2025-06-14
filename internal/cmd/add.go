package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/skewb1k/upfile/internal/service"
	"github.com/spf13/cobra"
)

func addCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <path>...",
		Short: "Start tracking specified file(s)",
		Args:  cobra.MinimumNArgs(1),
		// TODO: exclude tracked entries from completions
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				path, err := filepath.Abs(arg)
				if err != nil {
					return fmt.Errorf("failed to get absolute path to file '%s': %w", arg, err)
				}

				if err := service.Add(cmd.Context(), getIndexFsProvider(), path); err != nil {
					return fmt.Errorf("cannot add '%s': %w", path, err)
				}
			}

			return nil
		},
	}

	return cmd
}
