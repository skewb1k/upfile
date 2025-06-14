package cmd

import (
	"fmt"

	"github.com/skewb1k/upfile/internal/service"
	"github.com/spf13/cobra"
)

func renameCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rename <old> <new>",
		Short: "Rename all entries of file",
		Args:  cobra.ExactArgs(2),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			switch len(args) {
			case 0:
				return completeFname(cmd, args, toComplete)
			case 1:
				return nil, cobra.ShellCompDirectiveDefault
			default:
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			oldName, newName := args[0], args[1]
			if err := service.Rename(
				cmd.Context(),
				getStore(),
				oldName,
				newName,
			); err != nil {
				return fmt.Errorf("cannot rename to '%s': %w", newName, err)
			}

			return nil
		},
	}

	return cmd
}
