package common

// Image is a template for running a docker container
type Image struct {
	// Name is the docker image to base the container off of
	Name string `yaml:"image"`

	// Entrypoint is the entrypoint in the container
	Entrypoint []string

	// Envionrment is the environment vars to set on the container
	Environment []string

	// Hostname is the host name to set for the container
	Hostname string

	// Args are cli arguments to pass to the image
	Args []string `yaml:"command"`

	// Type is the container type, often service, batch, etc...
	Type string `yaml:"-"`

	// Volumes are volumes on the same engine
	Volumes []string

	// WorkingDir is the working directory to set for the container
	WorkingDir string `yaml:"working_dir"`

	// Pull tells the engine to always pull the Docker image
	Pull bool

	// Give extended privileges to this container, e.g. Docker-in-Docker
	Privileged bool

	// NetworkMode is the network mode for the container
	NetworkMode string `yaml:"net"`

	// UserData is user defined configuration data
	UserData map[string]interface{} `yaml:"user_data,inline"`
}
