package filename_test

import (
	"testing"

	"github.com/skewb1k/upfile/pkg/filename"
	"github.com/stretchr/testify/assert"
)

func TestExtract(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected string
		ok       bool
	}{
		{"/home/user/.prettierrc", ".prettierrc", true},
		{"relative/path/to/file.txt", "file.txt", true},
		{"file.txt", "file.txt", true},

		{"/", "", false},
		{".", "", false},
		{"", "", false},
		{"../", "", false},
		{"..", "", false},

		{"dir/../../../file", "file", true},
		{"a/b/../../..", "", false},

		{"./file", "file", true},
		{"./../file", "file", true},
		{"/tmp/../tmp/file", "file", true},
		{"dir/..", "", false},
		{"dir/../", "", false},

		{"some//nested///file", "file", true},
		{"/a//b/../../etc", "etc", true},

		{"./.env", ".env", true},
		{"..hidden", "..hidden", true},
		{"/dir/..hidden", "..hidden", true},

		{"folder/", "folder", true},
		{"folder////", "folder", true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			result, ok := filename.Extract(tt.input)
			assert.Equal(t, tt.expected, result, "unexpected result for input %q", tt.input)
			assert.Equal(t, tt.ok, ok, "unexpected ok for input %q", tt.input)
		})
	}
}
