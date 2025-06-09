package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	indexFs "github.com/skewb1k/upfile/internal/index/fs"
	"github.com/skewb1k/upfile/internal/service"
	userfileFs "github.com/skewb1k/upfile/internal/userfile/fs"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func getBaseDir() string {
	if dir := os.Getenv("UPFILE_DIR"); dir != "" {
		return dir
	}

	if xdgData := os.Getenv("XDG_DATA_HOME"); xdgData != "" {
		return filepath.Join(xdgData, name)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("failed to get current user's home dir: %s", err))
	}

	return filepath.Join(home, ".local", "share", name)
}

func wrap(f func(
	cmd *cobra.Command,
	s *service.Service,
	args []string,
) error,
) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		err := f(cmd, service.New(indexFs.New(getBaseDir()), userfileFs.New()), args)
		if errors.Is(err, service.ErrCancelled) {
			os.Exit(1)
		}

		return err
	}
}

func mustFprintf(w io.Writer, format string, a ...any) {
	if _, err := fmt.Fprintf(w, format, a...); err != nil {
		panic(err)
	}
}

// nolint: gochecknoglobals
var (
	green  = color.New(color.FgGreen).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
)

func statusAsString(status service.EntryStatus) string {
	switch status {
	case service.EntryStatusModified:
		return yellow("Modified")
	case service.EntryStatusUpToDate:
		return green("Up-to-date")
	case service.EntryStatusDeleted:
		return red("Deleted")
	default:
		panic("UNEXPECTED")
	}
}

func completeFname(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) ([]string, cobra.ShellCompDirective) {
	if len(args) >= 1 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	files, err := service.New(
		indexFs.New(getBaseDir()), userfileFs.New(),
	).ListTrackedFilenames(cmd.Context())
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	return files, cobra.ShellCompDirectiveNoFileComp
}

func doc(s string) string {
	return s[1:]
}
