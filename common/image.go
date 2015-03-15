package common

type Image struct {
	Name        string            `yaml:"image"`
	Type        string            `yaml:"type"`
	Env         map[string]string `yaml:"environment"`
	Cmd         []string          `yaml:"command"`
	Entrypoint  []string          `yaml:"entrypoint"`
	Volumes     []string          `yaml:"volumes"`
	Labels      []string          `yaml:"labels"`
	WorkingDir  string            `yaml:"working_dir"`
	Hostname    string            `yaml:"hostname"`
	NetworkMode string            `yaml:"net"`
	Privileged  bool              `yaml:"privileged"`
	Pull        bool              `yaml:"pull"`

	// UserData is user defined configuration data
	UserData map[string]interface{} `yaml:"user_data,inline"`
}
