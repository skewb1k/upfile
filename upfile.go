// Upfile manages syncing files with an upstream version.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func UpfileDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	upfileDir := filepath.Join(homeDir, ".local", "state", "upfile")
	return upfileDir, nil
}

func pushCmd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no file given")
	}
	// TODO: support providing multiple paths.
	if len(args) > 1 {
		return fmt.Errorf("1 file expected")
	}
	srcPath := args[0]

	upfileDir, err := UpfileDir()
	if err != nil {
		return err
	}
	upstreamDir := filepath.Join(upfileDir, "upstream")
	if err := os.MkdirAll(upstreamDir, 0o700); err != nil {
		return err
	}

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	upstream := filepath.Base(srcPath)
	upstreamPath := filepath.Join(upstreamDir, upstream)
	upstreamFile, err := os.OpenFile(upstreamPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer upstreamFile.Close()

	// TODO: perform atomic copy.
	_, err = upstreamFile.ReadFrom(srcFile)
	return err
}

func pullCmd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no upstream name given")
	}
	// TODO: support providing multiple upstreams.
	if len(args) > 1 {
		return fmt.Errorf("1 upstream name expected")
	}
	upstream := args[0]
	if !filepath.IsLocal(upstream) {
		return fmt.Errorf("invalid upstream name %s", upstream)
	}

	upfileDir, err := UpfileDir()
	if err != nil {
		return err
	}
	upstreamDir := filepath.Join(upfileDir, "upstream")

	path := filepath.Join(upstreamDir, upstream)
	src, err := os.Open(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("upstream file %s not found", upstream)
		}
		return err
	}
	defer src.Close()

	dst, err := os.Create(upstream)
	if err != nil {
		return err
	}
	defer dst.Close()

	// TODO: perform atomic copy.
	_, err = dst.ReadFrom(src)
	return err
}

func listCmd(args []string) error {
	upfileDir, err := UpfileDir()
	if err != nil {
		return err
	}
	upstreamDir := filepath.Join(upfileDir, "upstream")
	upstreams, err := os.ReadDir(upstreamDir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}
	for _, upstream := range upstreams {
		if upstream.IsDir() {
			continue
		}
		fmt.Println(upstream.Name())
	}
	return nil
}

func showCmd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no upstream name given")
	}
	// TODO: support providing multiple upstreams.
	if len(args) > 1 {
		return fmt.Errorf("1 upstream name expected")
	}
	upstream := args[0]
	if !filepath.IsLocal(upstream) {
		return fmt.Errorf("invalid upstream name %s", upstream)
	}

	upfileDir, err := UpfileDir()
	if err != nil {
		return err
	}
	upstreamDir := filepath.Join(upfileDir, "upstream")

	path := filepath.Join(upstreamDir, upstream)
	src, err := os.Open(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("upstream file %s not found", upstream)
		}
		return err
	}
	defer src.Close()

	_, err = os.Stdout.ReadFrom(src)
	return err
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "usage: %s <command> [args ...]\n", os.Args[0])
	fmt.Fprint(flag.CommandLine.Output(), "\nAvailable commands:\n\n")
	fmt.Fprint(flag.CommandLine.Output(), "  push   push file to upstream\n")
	fmt.Fprint(flag.CommandLine.Output(), "  pull   pull file from upstream\n")
	fmt.Fprint(flag.CommandLine.Output(), "  list   list upstreams\n")
	fmt.Fprint(flag.CommandLine.Output(), "  show   show upstream file content\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	log.SetPrefix(os.Args[0] + ": ")
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		log.Print("no command given")
		usage()
	}

	command := flag.Arg(0)

	switch command {
	case "push":
		if err := pushCmd(flag.Args()[1:]); err != nil {
			log.Fatalf("push: %s", err)
		}
	case "pull":
		if err := pullCmd(flag.Args()[1:]); err != nil {
			log.Fatalf("pull: %s", err)
		}
	case "list":
		if err := listCmd(flag.Args()[1:]); err != nil {
			log.Fatalf("list: %s", err)
		}
	case "show":
		if err := showCmd(flag.Args()[1:]); err != nil {
			log.Fatalf("show: %s", err)
		}
	default:
		log.Printf("unknown command %s", command)
		usage()
	}
}
