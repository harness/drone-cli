package common

// A config structure is used to configure a build.
type Config struct {
	Init  Step
	Clone Step
	Build Step

	Compose map[string]Step
	Publish map[string]Step
	Deploy  map[string]Step
	Notify  map[string]Step

	Matrix map[string][]string
}
