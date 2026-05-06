// Upfile manages syncing files with an upstream version.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	push = flag.Bool("push", false, "push file to upstream")
	pull = flag.Bool("pull", false, "pull file from upstream")
	list = flag.Bool("list", false, "list upstreams")
)

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "usage: upfile [flags] [args ...]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	log.SetPrefix("upfile: ")
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	if *push && *pull {
		usage()
	}
	if *push {
		err := pushCmd(flag.Args())
		if err != nil {
			log.Fatalf("push: %s", err)
		}
		return
	}
	if *pull {
		err := pullCmd(flag.Args())
		if err != nil {
			log.Fatalf("pull: %s", err)
		}
		return
	}
	if *list {
		err := listCmd(flag.Args())
		if err != nil {
			log.Fatalf("list: %s", err)
		}
		return
	}
	log.Print("no flag given")
	usage()
}

func UpfileDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	upfileDir := filepath.Join(homeDir, ".local", "state", "upfile")
	return upfileDir, nil
}
