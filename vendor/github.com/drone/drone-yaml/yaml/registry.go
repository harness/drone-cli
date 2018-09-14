package yaml

import "errors"

type (
	// Registry is a resource that provides encrypted
	// registry credentials and pointers to external
	// registry credentials (e.g. from vault).
	Registry struct {
		Kind string `json:"kind,omitempty"`
		Type string `json:"type,omitempty"`

		Data map[string]string `json:"data,omitempty"`
	}
)

// Validate returns an error if the registry is invalid.
func (r *Registry) Validate() error {
	if len(r.Data) == 0 {
		return errors.New("yaml: invalid registry resource")
	}
	return nil
}
