package main

import (
	"os"

	"github.com/skewb1k/upfile/internal/cmd"
)

func main() {
	os.Exit(cmd.Main(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
