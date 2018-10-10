package transform

import "github.com/drone/drone-runtime/engine"

// WithSecrets is a transform function that adds a set
// of global secrets to the container.
func WithSecrets(secrets map[string]string) func(*engine.Spec) {
	return func(spec *engine.Spec) {
		for key, value := range secrets {
			spec.Secrets = append(spec.Secrets,
				&engine.Secret{
					Name: key,
					Data: value,
				},
			)
		}
	}
}

// SecretFunc is a callback function used to request
// named secret, required by a pipeline step.
type SecretFunc func(string) *engine.Secret

// WithSecretFunc is a transform function that resolves
// all named secrets through a callback function, and
// adds the secrets to the specification.
func WithSecretFunc(f SecretFunc) func(*engine.Spec) {
	return func(spec *engine.Spec) {
		// first we get a unique list of all secrets
		// used by the specification.
		set := map[string]struct{}{}
		for _, step := range spec.Steps {
			// if we know the step is not going to run,
			// we can ignore any secrets that it requires.
			if step.RunPolicy == engine.RunNever {
				continue
			}
			for _, v := range step.Secrets {
				set[v.Name] = struct{}{}
			}
		}

		// next we use the callback function to
		// get the value for each secret, and append
		// to the specification.
		for name := range set {
			secret := f(name)
			if secret != nil {
				spec.Secrets = append(spec.Secrets, secret)
			}
		}
	}
}
