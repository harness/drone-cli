package yaml

import "errors"

type (
	// Secret is a resource that provides encrypted data
	// and pointers to external data (i.e. from vault).
	Secret struct {
		Kind string `json:"kind,omitempty"`
		Type string `json:"type,omitempty"`

		Data map[string]string `json:"data,omitempty"`
	}
)

// GetKind returns the resource kind.
func (s *Secret) GetKind() string { return s.Kind }

// Validate returns an error if the secret is invalid.
func (s *Secret) Validate() error {
	if len(s.Data) == 0 {
		return errors.New("yaml: invalid secret resource")
	}
	return nil
}
