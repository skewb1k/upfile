package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/skewb1k/upfile/internal/commands"
	"github.com/spf13/cobra"
)

// Populated by goreleaser during build
var version = "unknown"

func errorHandler(w io.Writer, _ fang.Styles, err error) {
	_, _ = fmt.Fprintln(w, "upfile:", err)
}

func main() {
	root := &cobra.Command{
		Use:     "upfile",
		Version: version,
		Short: `
Track and sync files across projects

Support tool on GitHub: https://github.com/skewb1k/upfile
`[1:],
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
	}

	root.AddCommand(
		commands.Add(),
		commands.Remove(),
		commands.Diff(),
		commands.Show(),
		commands.List(),
		commands.Status(),
		commands.Pull(),
		commands.Push(),
		commands.Sync(),
		commands.Drop(),
		commands.Rename(),
	)

	if err := fang.Execute(
		context.Background(),
		root,
		fang.WithVersion(root.Version),
		fang.WithErrorHandler(errorHandler),
		fang.WithColorSchemeFunc(fang.AnsiColorScheme),
	); err != nil {
		os.Exit(1)
	}
}
