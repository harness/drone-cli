package transform

import "github.com/drone/drone-runtime/engine"

// WithNetworks is a transform function that attaches a
// list of user-defined Docker networks to each step.
func WithNetworks(networks []string) func(*engine.Spec) {
	return func(spec *engine.Spec) {
		for _, step := range spec.Steps {
			step.Docker.Networks = append(
				step.Docker.Networks, networks...)
		}
	}
}
