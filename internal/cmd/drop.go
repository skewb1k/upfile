package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/skewb1k/upfile/internal/service"

	"github.com/spf13/cobra"
)

func dropCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "drop <filename>",
		Args:  cobra.ExactArgs(1),
		Short: "Remove tracked file upstream and entries",
		Long: doc(`
Permanently removes a file from UpFile tracking:

- Deletes the upstream version
- Removes all tracked entries from the index

Note: This does NOT delete any actual files from user-space filesystem.

Use with caution. You will be prompted to confirm removal unless --yes is specified.
`),

		RunE: wrap(func(cmd *cobra.Command, s *service.Service, args []string) error {
			fname := filepath.Base(args[0])

			confirm := func(entries []string) bool {
				if yes {
					return true
				}

				fmt.Println("The following files will be untracked and removed from UpFile:")
				for _, e := range entries {
					fmt.Println(" -", e)
				}
				fmt.Print("\nProceed? [y/N]: ")

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

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Automatic 'yes' to prompts")
	cmd.ValidArgsFunction = completeFname

	return cmd
}
