package validfname_test

import (
	"testing"

	"github.com/skewb1k/upfile/pkg/validfname"
	"github.com/stretchr/testify/assert"
)

func TestValidateFilename(t *testing.T) {
	t.Parallel()

	cases := []struct {
		filename string
		valid    bool
	}{
		{"", false},
		{".", false},
		{"..", false},
		{"file/name", false},
		{"bad\x00name", false},

		{"file", true},
		{"file.txt", true},
		{"some file.txt", true},
		{"file-name_1", true},
	}

	for _, tc := range cases {
		t.Run(tc.filename, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.valid, validfname.ValidateFilename(tc.filename))
		})
	}
}
