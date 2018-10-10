package transform

import "github.com/drone/drone-runtime/engine"

// WithEnviron is a transform function that adds a set
// of environment variables to each container.
func WithEnviron(envs map[string]string) func(*engine.Spec) {
	return func(spec *engine.Spec) {
		for key, value := range envs {
			for _, step := range spec.Steps {
				step.Envs[key] = value
			}
		}
	}
}
