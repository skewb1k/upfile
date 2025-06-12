package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/skewb1k/upfile/internal/service"
	"github.com/skewb1k/upfile/internal/store"
	"github.com/spf13/cobra"
)

func showCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "show <filename>",
		Short:             "Show upstream version of file",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completeFname,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Show(cmd.Context(), service.New(store.New(getBaseDir())), cmd.OutOrStdout(), args[0])
		},
	}

	return cmd
}

func Show(
	ctx context.Context,
	s *service.Service,
	out io.Writer,
	fname string,
) error {
	content, err := s.CatLatest(ctx, fname)
	if err != nil {
		return fmt.Errorf("cat: %w", err)
	}

	mustFprintf(out, "%s", content)

	return nil
}
