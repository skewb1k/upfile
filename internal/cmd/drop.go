package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/skewb1k/upfile/internal/entries"
	"github.com/skewb1k/upfile/internal/upstreams"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			baseDir := getBaseDir()
			upstreamsProvider := upstreams.NewProvider(baseDir)
			entriesProvider := entries.NewProvider(baseDir)

			fname := args[0]

			// TODO: collect errors
			e, err := entriesProvider.GetEntriesByFilename(cmd.Context(), fname)
			if err != nil {
				if errors.Is(err, entries.ErrInvalidFilename) {
					return ErrNotTracked
				}

				return fmt.Errorf("get entries by filename: %w", err)
			}

			// FIXME:
			if len(e) == 0 {
				return ErrNoEntries
			}

			if !yes && !ask(cmd.InOrStdin(), e, false, "The following entries will be untracked:") {
				os.Exit(1)
				return nil
			}

			if err := upstreamsProvider.DeleteUpstream(cmd.Context(), fname); err != nil {
				if errors.Is(err, upstreams.ErrNotFound) {
					return ErrNotTracked
				}

				return fmt.Errorf("delete upstream: %w", err)
			}

			for _, entry := range e {
				if err := entriesProvider.DeleteEntry(cmd.Context(), fname, entry); err != nil {
					return fmt.Errorf("delete entry: %w", err)
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Automatic 'yes' to prompts")
	cmd.ValidArgsFunction = completeFname

	return cmd
}
