package transform

import "github.com/drone/drone-runtime/engine"

// WithLables is a transform function that adds a set
// of labels to each resource.
func WithLables(labels map[string]string) func(*engine.Spec) {
	return func(spec *engine.Spec) {
		for k, v := range labels {
			spec.Metadata.Labels[k] = v
		}
		for _, resource := range spec.Docker.Volumes {
			for k, v := range labels {
				resource.Metadata.Labels[k] = v
			}
		}
		for _, resource := range spec.Steps {
			for k, v := range labels {
				resource.Metadata.Labels[k] = v
			}
		}
	}
}
