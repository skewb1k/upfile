package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

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
			fname := filepath.Base(args[0])

			confirm := func(entries []string) bool {
				if yes {
					return true
				}

				fmt.Println("The following tracked files will be updated:")
				for _, e := range entries {
					fmt.Println(" -", e)
				}

				fmt.Print("\nProceed? [Y/n]: ")

				var input string
				_, _ = fmt.Fscanln(cmd.InOrStdin(), &input)

				input = strings.ToLower(strings.TrimSpace(input))
				return input == "" || input == "y"
			}

			if err := s.Sync(cmd.Context(), fname, confirm); err != nil {
				if errors.Is(err, service.ErrUpToDate) {
					cmd.Println("Everything up-to-date")
					return nil
				}

				return err //nolint: wrapcheck
			}

			return nil
		}),
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Automatic 'yes' to prompts")
	cmd.ValidArgsFunction = completeFname

	// cmd.Flags().BoolVarP(&patch, "patch", "p", false, "Interactively apply changes per entry")

	return cmd
}
