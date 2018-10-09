package yaml

import "errors"

type (
	// Signature is a resource that provides an hmac
	// signature of combined resources. This signature
	// can be used to validate authenticity and prevent
	// tampering.
	Signature struct {
		Version string `json:"version,omitempty"`
		Kind    string `json:"kind"`
		Hmac    string `json:"hmac"`
	}
)

// GetVersion returns the resource version.
func (s *Signature) GetVersion() string { return s.Version }

// GetKind returns the resource kind.
func (s *Signature) GetKind() string { return s.Kind }

// Validate returns an error if the signature is invalid.
func (s Signature) Validate() error {
	if s.Hmac == "" {
		return errors.New("yaml: invalid signature. missing hash")
	}
	return nil
}
