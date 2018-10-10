package runtime

import "github.com/drone/drone-runtime/engine"

// Option configures a Runtime option.
type Option func(*Runtime)

// WithEngine sets the Runtime engine.
func WithEngine(e engine.Engine) Option {
	return func(r *Runtime) {
		r.engine = e
	}
}

// WithConfig sets the Runtime configuration.
func WithConfig(c *engine.Spec) Option {
	return func(r *Runtime) {
		r.config = c
	}
}

// WithHooks sets the Runtime tracer.
func WithHooks(h *Hook) Option {
	return func(r *Runtime) {
		if h != nil {
			r.hook = h
		}
	}
}
