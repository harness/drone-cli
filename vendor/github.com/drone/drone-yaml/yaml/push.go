package yaml

type (
	// Push configures a Docker push.
	Push struct {
		Image string `json:"image,omitempty"`
	}

	// push is a tempoary type used to unmarshal
	// the Push struct when long format is used.
	push struct {
		Image string `json:"image,omitempty"`
	}
)

// UnmarshalYAML implements yaml unmarshalling.
func (p *Push) UnmarshalYAML(unmarshal func(interface{}) error) error {
	d := new(push)
	err := unmarshal(&d.Image)
	if err != nil {
		err = unmarshal(d)
	}
	p.Image = d.Image
	return err
}
