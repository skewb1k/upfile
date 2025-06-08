package safejoin_test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/skewb1k/upfile/pkg/safejoin"
	"github.com/stretchr/testify/assert"
)

func TestSafeJoinFilename(t *testing.T) {
	t.Parallel()

	tests := []struct {
		filename string
		wantErr  bool
	}{
		// valid filenames
		{".prettierrc", false},
		{"file.txt", false},
		{"normal-name", false},
		{"UPPER_and_123.ext", false},
		{"file.name.with.dots", false},
		{"trailing-space ", false},
		{" space-at-start", false},
		{"~", false},
		{"'name'", false},
		{"'name with space'", false},

		// empty or reserved
		{"", true},
		{".", true},
		{"/", true},
		{"..", true},
		{".." + string(filepath.Separator) + "file", true},

		// nested paths
		{"a" + string(filepath.Separator) + "b", true},
		{"nested/dir/file", true},
		{"a/b/../c", true},

		// absolute paths
		{"/etc/passwd", true},
		{"/absolute/path", true},

		// Simulate a Windows-style absolute path
		{func() string {
			if runtime.GOOS == "windows" {
				return `C:\Windows\System32`
			}
			return "C:/Windows/System32"
		}(), true},
		{string(filepath.Separator) + "rooted", true},

		{"\x00", false},
	}

	for _, tc := range tests {
		t.Run(tc.filename, func(t *testing.T) {
			t.Parallel()

			res, err := safejoin.SafeJoinFilename("/home/user/", tc.filename)
			if tc.wantErr {
				assert.Error(t, err, res)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
