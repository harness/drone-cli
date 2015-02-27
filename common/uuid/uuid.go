package uuid

import (
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"io"
)

// CreateUUID is a helper function that will
// create a random, unique identifier.
func CreateUUID() string {
	c := sha1.New()
	r := CreateRandom()
	io.WriteString(c, string(r))
	s := fmt.Sprintf("%x", c.Sum(nil))
	return s[0:10]
}

// CreateRandom creates a random block of bytes
// that we can use to generate unique identifiers.
func CreateRandom() []byte {
	k := make([]byte, sha1.BlockSize)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return nil
	}
	return k
}
