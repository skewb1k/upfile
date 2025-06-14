package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func promptDefaultYes(stdin io.Reader, stdout io.Writer) (bool, error) {
	return prompt(stdin, stdout, true)
}

func promptDefaultNo(stdin io.Reader, stdout io.Writer) (bool, error) {
	return prompt(stdin, stdout, false)
}

func prompt(stdin io.Reader, stdout io.Writer, defaultYes bool) (bool, error) {
	var proceedMsg string
	if defaultYes {
		proceedMsg = "\nProceed? [Y/n]: "
	} else {
		proceedMsg = "\nProceed? [y/N]: "
	}

	if _, err := fmt.Fprint(stdout, proceedMsg); err != nil {
		return false, fmt.Errorf("failed to print proceed message: %w", err)
	}

	var input string
	_, _ = fmt.Fscanln(stdin, &input)

	input = strings.ToLower(strings.TrimSpace(input))

	if defaultYes && input == "" {
		return true, nil
	}

	return input == "y", nil
}

func MkdirAllWriteFile(path string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}

	if err := os.WriteFile(path, content, 0o600); err != nil {
		return err
	}

	return nil
}

func mustFmt(_ int, err error) {
	if err != nil {
		panic(err)
	}
}
