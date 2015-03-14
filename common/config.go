package common

import "strings"

// Config represents a repository build configuration.
type Config struct {
	Setup *Step
	Clone *Step
	Build *Step

	Compose map[string]*Step
	Publish map[string]*Step
	Deploy  map[string]*Step
	Notify  map[string]*Step

	Matrix Matrix
	Axis   Axis
}

// Matrix represents the build matrix.
type Matrix map[string][]string

// Axis represents a single permutation of entries
// from the build matrix.
type Axis map[string]string

// String returns a string representation of an Axis as
// a comma-separated list of environment variables.
func (a Axis) String() string {
	var envs []string
	for k, v := range a {
		envs = append(envs, k+"="+v)
	}
	return strings.Join(envs, " ")
}
