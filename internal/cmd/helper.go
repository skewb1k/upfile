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

// func wrap(f func(
// 	cmd *cobra.Command,
// 	s *service.Service,
// 	args []string,
// ) error,
// ) func(cmd *cobra.Command, args []string) error {
// 	return func(cmd *cobra.Command, args []string) error {
// 		err := f(cmd, service.New(indexFs.New(getBaseDir()), userfileFs.New()), args)
// 		if errors.Is(err, service.ErrCancelled) {
// 			os.Exit(1)
// 		}
//
// 		return err
// 	}
// }

// func mustFprintf(w io.Writer, format string, a ...any) {
// 	if _, err := fmt.Fprintf(w, format, a...); err != nil {
// 		panic(err)
// 	}
// }
//
// func mustFprintln(w io.Writer, a ...any) {
// 	if _, err := fmt.Fprintln(w, a...); err != nil {
// 		panic(err)
// 	}
// }

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

func ask(stdin io.Reader, list []string, defaultYes bool, msg string) bool {
	fmt.Println(msg)
	for _, e := range list {
		fmt.Println(" -", e)
	}
	if defaultYes {
		fmt.Print("\nProceed? [Y/n]: ")
	} else {
		fmt.Print("\nProceed? [y/N]: ")
	}

	var input string
	_, _ = fmt.Fscanln(stdin, &input)

	input = strings.ToLower(strings.TrimSpace(input))

	if defaultYes && input == "" {
		return true
	}

	return input == "y"
}

func WriteFile(path string, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create parent dirs: %w", err)
	}

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
