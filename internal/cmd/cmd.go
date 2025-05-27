package cmd

import (
	"context"
	"errors"
	"io"
	"os/exec"

	cc "github.com/ivanpirog/coloredcobra"
	"github.com/spf13/cobra"
)

func Main(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) int {
	rootCmd := &cobra.Command{
		Use: "upfile",
	}

	// TODO: adjust colors
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
	rootCmd.AddCommand(link())
	rootCmd.AddCommand(diff())

	ctx := context.Background()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}

		return 1
	}

	return 0
}
