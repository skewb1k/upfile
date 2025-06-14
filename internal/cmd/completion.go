package cmd

import (
	"github.com/spf13/cobra"
)

func completeFnames(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) ([]string, cobra.ShellCompDirective) {
	files, err := getStore().GetFilenames(cmd.Context())
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	// Create a map for quick lookup of existing arguments
	argMap := make(map[string]struct{})
	for _, arg := range args {
		argMap[arg] = struct{}{}
	}

	// Filter out files that are already present in args
	var filteredFiles []string
	for _, file := range files {
		if _, found := argMap[file]; !found {
			filteredFiles = append(filteredFiles, file)
		}
	}

	return filteredFiles, cobra.ShellCompDirectiveNoFileComp
}

func completeFname(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) ([]string, cobra.ShellCompDirective) {
	if len(args) >= 1 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	files, err := getStore().GetFilenames(cmd.Context())
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	return files, cobra.ShellCompDirectiveNoFileComp
}
