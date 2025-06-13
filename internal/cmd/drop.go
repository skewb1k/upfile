package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/skewb1k/upfile/internal/store"
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
		ValidArgsFunction: completeFname,
		RunE: withStore(func(cmd *cobra.Command, s *store.Store, args []string) error {
			fname := args[0]

			entries, err := s.GetEntriesByFilename(cmd.Context(), fname)
			if err != nil {
				return fmt.Errorf("get entries by filename: %w", err)
			}

			if len(entries) != 0 && !yes {
				cmd.Println("The following entries will be untracked:")

				for _, e := range entries {
					cmd.Println(" -", e)
				}

				ok, err := askDefaultNo(cmd.InOrStdin(), cmd.OutOrStdout())
				if err != nil {
					return err
				}

				if !ok {
					os.Exit(1)
				}
			}

			if err := s.DeleteUpstream(cmd.Context(), fname); err != nil {
				if errors.Is(err, store.ErrNotFound) {
					return ErrNotTracked
				}

				return fmt.Errorf("delete upstream: %w", err)
			}

			for _, entry := range entries {
				if err := s.DeleteEntry(cmd.Context(), fname, entry); err != nil {
					return fmt.Errorf("delete entry: %w", err)
				}
			}

			return nil
		}),
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Automatic 'yes' to prompts")

	return cmd
}
