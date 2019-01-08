package bitbucket

import (
	"path"
	"strings"
)

type (
	// Config defines the pipeline configuration.
	Config struct {
		// Image specifies the Docker image with
		// which we run your builds.
		Image string

		// Clone defines the depth of Git clones
		// for all pipelines.
		Clone struct {
			Depth int
		}

		// Pipeline defines the pipeline configuration
		// which includes a list of all steps for default,
		// tag, and branch-specific execution.
		Pipelines struct {
			Default  Stage
			Tags     map[string]Stage
			Branches map[string]Stage
		}

		Definitions struct {
			Services map[string]*Step
			Caches   map[string]string
		}
	}

	// Stage contains a list of steps executed
	// for a specific branch or tag.
	Stage struct {
		Name  string
		Steps []*Step
	}

	// Step defines a build execution unit.
	Step struct {
		// Name of the pipeline step.
		Name string

		// Image specifies the Docker image with
		// which we run your builds.
		Image string

		// Script contains the list of bash commands
		// that are executed in sequence.
		Script []string

		// Variables provides environment variables
		// passed to the container at runtime.
		Variables map[string]string

		// Artifacts defines files that are to be
		// snapshotted and shared with the subsequent
		// step. This is not used, because Drone uses
		// a shared volume to share artifacts.
		Artifacts []string
	}
)

// Pipeline returns the pipeline stage that best matches the branch
// and ref. If there is no matching pipeline specific to the branch
// or tag, the default pipeline is returned.
func (c *Config) Pipeline(ref string) Stage {
	// match pipeline by tag name
	tag := strings.TrimPrefix(ref, "refs/tags/")
	for pattern, pipeline := range c.Pipelines.Tags {
		if ok, _ := path.Match(pattern, tag); ok {
			return pipeline
		}
	}
	// match pipeline by branch name
	branch := strings.TrimPrefix(ref, "refs/heads/")
	for pattern, pipeline := range c.Pipelines.Branches {
		if ok, _ := path.Match(pattern, branch); ok {
			return pipeline
		}
	}
	// use default
	return c.Pipelines.Default
}

// UnmarshalYAML implements custom parsing for the stage section of the yaml
// to cleanup the structure a bit.
func (s *Stage) UnmarshalYAML(unmarshal func(interface{}) error) error {
	in := []struct {
		Step *Step
	}{}
	err := unmarshal(&in)
	if err != nil {
		return err
	}
	for _, step := range in {
		s.Steps = append(s.Steps, step.Step)
	}
	return nil
}
