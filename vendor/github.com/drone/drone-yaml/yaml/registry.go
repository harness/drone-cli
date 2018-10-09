package yaml

import "errors"

type (
	// Registry is a resource that provides encrypted
	// registry credentials and pointers to external
	// registry credentials (e.g. from vault).
	Registry struct {
		Version string `json:"version,omitempt"`
		Kind    string `json:"kind,omitempty"`
		Type    string `json:"type,omitempty"`

		Data map[string]string `json:"data,omitempty"`
	}
)

// GetVersion returns the resource version.
func (r *Registry) GetVersion() string { return r.Version }

// GetKind returns the resource kind.
func (r *Registry) GetKind() string { return r.Kind }

// Validate returns an error if the registry is invalid.
func (r *Registry) Validate() error {
	if len(r.Data) == 0 {
		return errors.New("yaml: invalid registry resource")
	}
	return nil
}
