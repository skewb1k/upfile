package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"upfile/internal/service"

	"github.com/spf13/cobra"
)

func drop() *cobra.Command {
	return &cobra.Command{
		Use:   "drop <filename>",
		Short: "Drop all entries and upstream version of file",
		Args:  cobra.ExactArgs(1),
		RunE: wrap(func(cmd *cobra.Command, s *service.Service, args []string) error {
			fname := filepath.Base(args[0])

			confirm := func(entries []string) bool {
				fmt.Println("The following tracked paths will be deleted:")
				for _, e := range entries {
					fmt.Println(" -", e)
				}
				fmt.Print("Proceed? [y/N]: ")

				var input string
				_, _ = fmt.Fscanln(cmd.InOrStdin(), &input)

				return strings.ToLower(strings.TrimSpace(input)) == "y"
			}

			if err := s.Drop(cmd.Context(), fname, confirm); err != nil {
				return err //nolint: wrapcheck
			}

			return nil
		}),
	}
}
