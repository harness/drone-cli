package common

// A step represents a step in the build process, including
// the execution environment and parameters.
type Step struct {
	Name        string
	Image       string
	Environment []string
	Volumes     []string
	Hostname    string
	Privileged  bool
	Net         string

	// Config represents the unique configuration details
	// for each plugin.
	Config map[string]interface{} `yaml:"config,inline"`
}
