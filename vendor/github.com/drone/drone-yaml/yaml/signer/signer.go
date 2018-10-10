package signer

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/drone/drone-yaml/yaml"

	goyaml "gopkg.in/yaml.v2"
)

// ErrInvalidKey is returned when the key is missing or
// is less than 32-bytes.
var ErrInvalidKey = errors.New("signer: key must be 32-bytes")

// Key represents 32-byte signature.
type Key []byte

// KeyString is a helper function that returns a Key
// from a string.
func KeyString(s string) Key {
	return []byte(s)
}

// Sign calculates and returns the hmac signature of the
// parsed yaml file.
func Sign(data []byte, key Key) (string, error) {
	res, err := yaml.ParseRawBytes(data)
	if err != nil {
		return "", err
	}
	hmac, err := sign(res, key)
	return hex.EncodeToString(hmac), err
}

// SignUpdate calculates the hmac signature of the parsed
// yaml file and adds a signature resource. If a signature
// resource already exists, it is replaced.
func SignUpdate(data []byte, key Key) ([]byte, error) {
	res, err := yaml.ParseRawBytes(data)
	if err != nil {
		return nil, err
	}
	hmac, err := sign(res, key)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	for _, r := range res {
		if r.Kind != yaml.KindSignature {
			buf.WriteString("---")
			buf.WriteByte('\n')
			buf.Write(r.Data)
		}
	}

	buf.WriteString("---")
	buf.WriteByte('\n')
	buf.WriteString("kind: signature")
	buf.WriteByte('\n')
	buf.WriteString("hmac: " + hex.EncodeToString(hmac))
	buf.WriteByte('\n')
	buf.WriteByte('\n')
	buf.WriteString("...")
	buf.WriteByte('\n')
	return buf.Bytes(), nil
}

// Verify returns true if the signature of the parsed
// yaml file can be verified.
func Verify(data []byte, key Key) (bool, error) {
	res, err := yaml.ParseRawBytes(data)
	if err != nil {
		return false, err
	}
	mac1, err := extract(res)
	if err != nil {
		return false, nil
	}
	mac2, err := sign(res, key)
	if err != nil {
		return false, err
	}
	return hmac.Equal(mac1, mac2), nil
}

// WriteTo writes the signature to the yaml file. If the
// signature already exists it is removed, and the new
// signature is appended to the end of the document.
func WriteTo(data []byte, hmac string) ([]byte, error) {
	res, err := yaml.ParseRawBytes(data)
	return upsert(res, hmac), err
}

// helper function extracts the hex-encoded signature
// resource from the parsed resource list.
func extract(res []*yaml.RawResource) ([]byte, error) {
	for _, r := range res {
		if r.Kind == yaml.KindSignature {
			out := new(yaml.Signature)
			err := goyaml.Unmarshal(r.Data, out)
			if err != nil {
				return nil, err
			}
			return hex.DecodeString(out.Hmac)
		}
	}
	return nil, errors.New("yaml: missing signature")
}

// helper function generates a hex-encoded signature
// based on the parsed resource list.
func sign(resources []*yaml.RawResource, key Key) ([]byte, error) {
	if len(key) < 32 {
		return nil, ErrInvalidKey
	}
	h := hmac.New(sha256.New, key)
	for _, r := range resources {
		if r.Kind != yaml.KindSignature {
			h.Write(r.Data)
		}
	}
	return h.Sum(nil), nil
}

// helper function inserts or updates the hmac signature
// into the yaml document, and returns an updated copy.
func upsert(res []*yaml.RawResource, hmac string) []byte {
	var buf bytes.Buffer
	for _, r := range res {
		if r.Kind != yaml.KindSignature {
			buf.WriteString("---")
			buf.WriteByte('\n')
			buf.Write(r.Data)
		}
	}
	buf.WriteString("---")
	buf.WriteByte('\n')
	buf.WriteString("kind: signature")
	buf.WriteByte('\n')
	buf.WriteString("hmac: " + hmac)
	buf.WriteByte('\n')
	buf.WriteByte('\n')
	buf.WriteString("...")
	buf.WriteByte('\n')
	return buf.Bytes()
}
