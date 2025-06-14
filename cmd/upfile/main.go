package main

import (
	"os"

	"github.com/skewb1k/upfile/internal/cmd"
)

// Populated by goreleaser during build
var version = "unknown"

func main() {
	os.Exit(cmd.Main(
		version,
		os.Args[1:],
		os.Stdin,
		os.Stdout,
		os.Stderr,
	))
}
