package cmd

import (
	"errors"
	"io"
	"os/exec"

	cc "github.com/ivanpirog/coloredcobra"

	"github.com/spf13/cobra"
)

const Name = "upfile"

func Main(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	rootCmd := &cobra.Command{
		Use: Name,
		Long: `
Sync files across multiple projects
skewb1k <skewb1kunix@gmail.com>
Source: https://github.com/skewb1k/upfile
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

	rootCmd.AddCommand(version())
	rootCmd.AddCommand(add())
	rootCmd.AddCommand(remove())
	rootCmd.AddCommand(diff())
	rootCmd.AddCommand(show())
	rootCmd.AddCommand(list())
	rootCmd.AddCommand(status())
	rootCmd.AddCommand(pull())
	rootCmd.AddCommand(push())
	rootCmd.AddCommand(sync())
	rootCmd.AddCommand(drop())

	if err := rootCmd.Execute(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}

		return 1
	}

	return 0
}
