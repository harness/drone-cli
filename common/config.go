package common

// Config represents a repository build configuration.
type Config struct {
	Setup *Step
	Clone *Step
	Build *Step

	Compose map[string]*Step
	Publish map[string]*Step
	Deploy  map[string]*Step
	Notify  map[string]*Step

	Matrix map[string][]string
	Axis   map[string]string
}
