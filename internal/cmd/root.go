package cmd

import (
	"errors"
	"io"
	"os/exec"

	cc "github.com/ivanpirog/coloredcobra"

	"github.com/spf13/cobra"
)

func Main(
	version string,
	args []string,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) int {
	cmd := &cobra.Command{
		Use:     "upfile",
		Version: version + "\n",
		Short: `
Track and sync files across projects

Support tool on Github: https://github.com/skewb1k/upfile
`[1:],
		SilenceUsage:      true,
		SilenceErrors:     false,
		CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
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
