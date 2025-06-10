package sha256

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

const Size = 32

type SHA256 [Size]byte

var ErrInvalidLength = errors.New("invalid length")

func ConvertSlice(hash []byte) (SHA256, error) {
	var sha SHA256

	if len(hash) != Size {
		return SHA256{}, ErrInvalidLength
	}

	copy(sha[:], hash)
	return sha, nil
}

func (h SHA256) EqualString(s string) bool {
	return h == FromString(s)
}

func (h SHA256) EqualBytes(s []byte) bool {
	return h == FromBytes(s)
}

func (h SHA256) String() string {
	return hex.EncodeToString(h[:])
}

func FromString(s string) SHA256 {
	return sha256.Sum256([]byte(s))
}

func FromBytes(s []byte) SHA256 {
	return sha256.Sum256(s)
}
