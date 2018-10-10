package transform

import "github.com/drone/drone-runtime/engine"

// WithAuths is a transform function that adds a set
// of global registry credentials to the container.
func WithAuths(auths []*engine.DockerAuth) func(*engine.Spec) {
	return func(spec *engine.Spec) {
		for _, auth := range auths {
			spec.Docker.Auths = append(spec.Docker.Auths, auth)
		}
	}
}

// AuthsFunc is a callback function used to request
// registry credentials to pull private images.
type AuthsFunc func() []*engine.DockerAuth

// WithAuthsFunc is a transform function that provides
// the sepcification with registry authentication
// credentials via a callback function.
func WithAuthsFunc(f AuthsFunc) func(*engine.Spec) {
	return func(spec *engine.Spec) {
		auths := f()
		if len(auths) != 0 {
			spec.Docker.Auths = append(spec.Docker.Auths, auths...)
		}
	}
}
