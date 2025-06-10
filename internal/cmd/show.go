package cmd

import (
	"fmt"

	"github.com/skewb1k/upfile/internal/upstreams"
	"github.com/spf13/cobra"
)

func showCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <filename>",
		Short: "Show upstream version of file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			baseDir := getBaseDir()
			upstreamsProvider := upstreams.NewProvider(baseDir)

			fname := args[0]

			upstream, err := upstreamsProvider.GetUpstream(cmd.Context(), fname)
			if err != nil {
				// if errors.Is(err, index.ErrNotFound) || errors.Is(err, index.ErrInvalidFilename) {
				// 	return ErrNotTracked
				// }
				//
				return fmt.Errorf("get upstream: %w", err)
			}

			cmd.Print(upstream.Content)

			return nil
		},
	}

	cmd.ValidArgsFunction = completeFname

	return cmd
}
