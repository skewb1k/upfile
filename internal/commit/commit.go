package commit

import (
	stdsha "crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/skewb1k/upfile/pkg/sha256"
)

func (c Commit) EncodeCommit(w io.Writer) {
	if err := json.NewEncoder(w).Encode(c); err != nil {
		panic(err)
	}
}

func DecodeCommit(r io.Reader) (*Commit, error) {
	var c Commit

	if err := json.NewDecoder(r).Decode(&c); err != nil {
		return nil, fmt.Errorf("decode commit: %w", err)
	}

	return &c, nil
}

type Commit struct {
	Filename    string
	ContentHash sha256.SHA256
	Parent      *sha256.SHA256
	Timestamp   time.Time
}

func (c Commit) HashCommit() sha256.SHA256 {
	h := stdsha.New()

	c.EncodeCommit(h)

	return sha256.FromBytes(h.Sum(nil))
}
