package yaml

import "errors"

type (
	// Signature is a resource that provides an hmac
	// signature of combined resources. This signature
	// can be used to validate authenticity and prevent
	// tampering.
	Signature struct {
		Kind string `json:"kind"`
		Hmac string `json:"hmac"`
	}
)

// Validate returns an error if the signature is invalid.
func (s Signature) Validate() error {
	if s.Hmac == "" {
		return errors.New("yaml: invalid signature. missing hash")
	}
	return nil
}
