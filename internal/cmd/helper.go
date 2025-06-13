package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/skewb1k/upfile/internal/store"
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

func withStore(f func(
	cmd *cobra.Command,
	s *store.Store,
	args []string,
) error,
) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return f(cmd, store.New(getBaseDir()), args)
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

	s := store.New(getBaseDir())
	files, err := s.GetFilenames(cmd.Context())
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	return files, cobra.ShellCompDirectiveNoFileComp
}

func doc(s string) string {
	return s[1:]
}

func askDefaultYes(stdin io.Reader, stdout io.Writer) (bool, error) {
	return ask(stdin, stdout, true)
}

func askDefaultNo(stdin io.Reader, stdout io.Writer) (bool, error) {
	return ask(stdin, stdout, false)
}

func ask(stdin io.Reader, stdout io.Writer, defaultYes bool) (bool, error) {
	var proceedMsg string
	if defaultYes {
		proceedMsg = "\nProceed? [Y/n]: "
	} else {
		proceedMsg = "\nProceed? [y/N]: "
	}

	if _, err := fmt.Fprint(stdout, proceedMsg); err != nil {
		return false, fmt.Errorf("failed to print proceed message: %w", err)
	}

	var input string
	_, _ = fmt.Fscanln(stdin, &input)

	input = strings.ToLower(strings.TrimSpace(input))

	if defaultYes && input == "" {
		return true, nil
	}

	return input == "y", nil
}

func MkdirAllWriteFile(path string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create parent dirs: %w", err)
	}

	if err := os.WriteFile(path, content, 0o600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
