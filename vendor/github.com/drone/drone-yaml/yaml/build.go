package yaml

type (
	// Build configures a Docker build.
	Build struct {
		Args       map[string]string `json:"args,omitempty"`
		CacheFrom  []string          `json:"cache_from,omitempty" yaml:"cache_from"`
		Context    string            `json:"context,omitempty"`
		Dockerfile string            `json:"dockerfile,omitempty"`
		Image      string            `json:"image,omitempty"`
		Labels     map[string]string `json:"labels,omitempty"`
	}

	// build is a tempoary type used to unmarshal
	// the Build struct when long format is used.
	build struct {
		Args       map[string]string
		CacheFrom  []string `yaml:"cache_from"`
		Context    string
		Dockerfile string
		Image      string
		Labels     map[string]string
	}
)

// UnmarshalYAML implements yaml unmarshalling.
func (b *Build) UnmarshalYAML(unmarshal func(interface{}) error) error {
	d := new(build)
	err := unmarshal(&d.Image)
	if err != nil {
		err = unmarshal(d)
	}
	b.Args = d.Args
	b.CacheFrom = d.CacheFrom
	b.Context = d.Context
	b.Dockerfile = d.Dockerfile
	b.Labels = d.Labels
	b.Image = d.Image
	return err
}
