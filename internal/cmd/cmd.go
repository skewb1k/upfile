package cmd

import (
	"errors"
	"io"
	"os/exec"

	cc "github.com/ivanpirog/coloredcobra"

	"github.com/spf13/cobra"
)

const name = "upfile"

// Populated by goreleaser during build
var version = "unknown"

func Main(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	rootCmd := &cobra.Command{
		Use: name,
		Short: `
Track and sync files across projects

Support project on Github: https://github.com/skewb1k/upfile
`[1:],
		SilenceUsage:  true,
		SilenceErrors: false,
	}
	rootCmd.SetErrPrefix(red("error:"))

	cc.Init(&cc.Config{
		RootCmd:         rootCmd,
		Headings:        cc.HiCyan + cc.Bold + cc.Underline,
		Commands:        cc.HiYellow + cc.Bold,
		Example:         cc.Italic,
		ExecName:        cc.Bold,
		Flags:           cc.Bold,
		NoExtraNewlines: true,
	})

	rootCmd.SetArgs(args)
	rootCmd.SetIn(stdin)
	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)

	rootCmd.AddCommand(versionCmd())
	rootCmd.AddCommand(addCmd())
	rootCmd.AddCommand(removeCmd())
	rootCmd.AddCommand(diffCmd())
	rootCmd.AddCommand(showCmd())
	rootCmd.AddCommand(listCmd())
	rootCmd.AddCommand(statusCmd())
	rootCmd.AddCommand(pullCmd())
	rootCmd.AddCommand(pushCmd())
	rootCmd.AddCommand(syncCmd())
	rootCmd.AddCommand(dropCmd())

	if err := rootCmd.Execute(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}

		return 1
	}

	return 0
}
