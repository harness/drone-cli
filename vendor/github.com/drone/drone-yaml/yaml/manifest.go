package yaml

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

// Resource enums.
const (
	KindCron      = "cron"
	KindPipeline  = "pipeline"
	KindRegistry  = "registry"
	KindSecret    = "secret"
	KindSignature = "signature"
)

type (
	// Manifest is a collection of Drone resources.
	Manifest struct {
		Resources []Resource
	}

	// Resource represents a Drone resource.
	Resource interface {
		// GetKind returns the resource kind.
		GetKind() string
	}

	// RawResource is a raw encoded resource with the
	// resource kind and type extracted.
	RawResource struct {
		Kind string
		Type string
		Data []byte `yaml:"-"`
	}

	resource struct {
		Kind string `json:"kind"`
		Type string `json:"type"`
	}
)

// UnmarshalJSON implement the json.Unmarshaler.
func (m *Manifest) UnmarshalJSON(b []byte) error {
	messages := []json.RawMessage{}
	err := json.Unmarshal(b, &messages)
	if err != nil {
		return err
	}
	for _, message := range messages {
		res := new(resource)
		err := json.Unmarshal(message, res)
		if err != nil {
			return err
		}
		var obj Resource
		switch res.Kind {
		case "cron":
			obj = new(Cron)
		case "secret":
			obj = new(Secret)
		case "signature":
			obj = new(Signature)
		case "registry":
			obj = new(Registry)
		default:
			obj = new(Pipeline)
		}
		err = json.Unmarshal(message, obj)
		if err != nil {
			return err
		}
		m.Resources = append(m.Resources, obj)
	}
	return nil
}

// MarshalJSON implement the json.Marshaler.
func (m *Manifest) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Resources)
}

// MarshalYAML implement the yaml.Marshaler.
func (m *Manifest) MarshalYAML() (interface{}, error) {
	return yaml.Marshal(m.Resources)
}
