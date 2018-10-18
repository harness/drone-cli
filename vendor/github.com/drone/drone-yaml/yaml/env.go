package yaml

type (
	// Variable represents an environment variable that
	// can be defined as a string literal or as a reference
	// to a secret.
	Variable struct {
		Value  string `json:"value,omitempty"`
		Secret string `json:"from_secret,omitempty" yaml:"from_secret"`
	}

	// variable is a tempoary type used to unmarshal
	// variables with references to secrets.
	variable struct {
		Value  string
		Secret string `yaml:"from_secret"`
	}
)

// UnmarshalYAML implements yaml unmarhsaling.
func (v *Variable) UnmarshalYAML(unmarshal func(interface{}) error) error {
	d := new(variable)
	err := unmarshal(&d.Value)
	if err != nil {
		err = unmarshal(d)
	}
	v.Value = d.Value
	v.Secret = d.Secret
	return err
}
