package yaml

type (

	// Port represents a network port in a single container.
	Port struct {
		Port     int    `json:"port,omitempty"`
		Host     int    `json:"host,omitempty"`
		Protocol string `json:"protocol,omitempty"`
	}

	port struct {
		Port     int
		Host     int
		Protocol string
	}
)

// UnmarshalYAML implements yaml unmarhsaling.
func (p *Port) UnmarshalYAML(unmarshal func(interface{}) error) error {
	out := new(port)
	err := unmarshal(&out.Port)
	if err != nil {
		err = unmarshal(&out)
	}
	p.Port = out.Port
	p.Host = out.Host
	p.Protocol = out.Protocol
	return err
}
