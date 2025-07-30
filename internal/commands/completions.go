package commands

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

func completeFnames(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) ([]string, cobra.ShellCompDirective) {
	files, err := getIndexFsProvider().GetFilenames(cmd.Context())
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	argMap := make(map[string]struct{})
	for _, arg := range args {
		argMap[arg] = struct{}{}
	}

	var filteredFiles []string
	for _, file := range files {
		if _, found := argMap[file]; !found {
			filteredFiles = append(filteredFiles, escapeBackslashes(file))
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

	files, err := getIndexFsProvider().GetFilenames(cmd.Context())
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	for i, file := range files {
		files[i] = escapeBackslashes(file)
	}

	return files, cobra.ShellCompDirectiveNoFileComp
}

func completeEntry(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	showHidden := len(toComplete) > 0 && strings.HasPrefix(filepath.Base(toComplete), ".")

	// Normalize and join input
	toComplete = filepath.Clean(toComplete)
	candidate := filepath.Join(cwd, toComplete)

	dir := candidate
	underDir := true
	if fi, err := os.Stat(candidate); err != nil || !fi.IsDir() {
		dir = filepath.Dir(candidate)
		underDir = false
	}

	indexProvider := getIndexFsProvider()
	files, err := indexProvider.GetFilenamesByEntry(cmd.Context(), dir)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	if underDir {
		for i, f := range files {
			files[i] = filepath.Join(toComplete, f)
		}
	}

	if entries, err := os.ReadDir(dir); err == nil {
		for _, entry := range entries {
			name := entry.Name()

			if entry.IsDir() {
				if showHidden || !strings.HasPrefix(name, ".") {
					var suggestion string
					if underDir {
						suggestion = filepath.Join(toComplete, name)
					} else {
						suggestion = name
					}

					files = append(files, suggestion+string(filepath.Separator))
				}
			}
		}
	}

	for i, f := range files {
		files[i] = escapeBackslashes(f)
	}

	sort.Strings(files)
	return files, cobra.ShellCompDirectiveNoSpace
}

func escapeBackslashes(s string) string {
	return strings.ReplaceAll(s, `\`, `\\`)
}
