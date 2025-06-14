package cmd

import (
	"errors"
	"io"
	"os/exec"

	cc "github.com/ivanpirog/coloredcobra"

	"github.com/spf13/cobra"
)

// Populated by goreleaser during build
var version = "unknown"

func Main(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	cmd := &cobra.Command{
		Use:     "upfile",
		Version: version + "\n",
		Short: `
Track and sync files across projects

Support project on Github: https://github.com/skewb1k/upfile
`[1:],
		SilenceUsage:  true,
		SilenceErrors: false,
	}

	cmd.SetVersionTemplate("{{.Version}}")
	cmd.SetErrPrefix("upfile:")

	cc.Init(&cc.Config{
		RootCmd:         cmd,
		Headings:        cc.HiCyan + cc.Bold + cc.Underline,
		Commands:        cc.HiYellow + cc.Bold,
		Example:         cc.Italic,
		ExecName:        cc.Bold,
		Flags:           cc.Bold,
		NoExtraNewlines: true,
	})

	cmd.SetArgs(args)
	cmd.SetIn(stdin)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)

	cmd.AddCommand(
		versionCmd(),
		addCmd(),
		removeCmd(),
		diffCmd(),
		showCmd(),
		listCmd(),
		statusCmd(),
		pullCmd(),
		pushCmd(),
		syncCmd(),
		dropCmd(),
		renameCmd(),
	)

	if err := cmd.Execute(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}

		return 1
	}

	return 0
}
