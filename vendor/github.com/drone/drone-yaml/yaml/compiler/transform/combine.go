package transform

import "github.com/drone/drone-runtime/engine"

// Combine is a transform function that combines
// one or many transform functions.
func Combine(fns ...func(*engine.Spec)) func(*engine.Spec) {
	return func(spec *engine.Spec) {
		for _, fn := range fns {
			fn(spec)
		}
	}
}
