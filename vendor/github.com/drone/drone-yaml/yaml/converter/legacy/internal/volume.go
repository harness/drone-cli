package yaml

import "strings"

// Volume represent a container volume.
type Volume struct {
	Source      string
	Destination string
	ReadOnly    bool
}

// UnmarshalYAML implements the Unmarshaller interface.
func (v *Volume) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var stringType string
	if err := unmarshal(&stringType); err != nil {
		return err
	}
	parts := strings.SplitN(stringType, ":", 3)
	switch {
	case len(parts) == 2:
		v.Source = parts[0]
		v.Destination = parts[1]
	case len(parts) == 3:
		v.Source = parts[0]
		v.Destination = parts[1]
		v.ReadOnly = parts[2] == "ro"
	}
	return nil
}
