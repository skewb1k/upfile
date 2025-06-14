package validfname

import "strings"

func ValidateFilename(fname string) bool {
	if fname == "" || fname == "." || fname == ".." {
		return false
	}

	if strings.ContainsRune(fname, '/') || strings.ContainsRune(fname, '\x00') {
		return false
	}

	return true
}
