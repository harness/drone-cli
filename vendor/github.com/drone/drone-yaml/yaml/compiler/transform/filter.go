package transform

import "github.com/drone/drone-runtime/engine"

// Include is a transform function that limits the
// pipeline execution to a whitelist of named steps.
func Include(names []string) func(*engine.Spec) {
	set := map[string]struct{}{}
	for _, name := range names {
		set[name] = struct{}{}
	}
	return func(spec *engine.Spec) {
		if len(names) == 0 {
			return
		}
		for _, step := range spec.Steps {
			if step.Metadata.Name == "clone" {
				continue
			}
			_, ok := set[step.Metadata.Name]
			if !ok {
				// if the step is not included in the
				// whitelist the run policy is set to never.
				step.RunPolicy = engine.RunNever
			}
		}
	}
}

// Exclude is a transform function that limits the
// pipeline execution to a whitelist of named steps.
func Exclude(names []string) func(*engine.Spec) {
	set := map[string]struct{}{}
	for _, name := range names {
		set[name] = struct{}{}
	}
	return func(spec *engine.Spec) {
		if len(names) == 0 {
			return
		}
		for _, step := range spec.Steps {
			if step.Metadata.Name == "clone" {
				continue
			}
			_, ok := set[step.Metadata.Name]
			if ok {
				// if the step is included in the blacklist
				// the run policy is set to never.
				step.RunPolicy = engine.RunNever
			}
		}
	}
}

// ResumeAt is a transform function that modifies the
// exuction to resume at a named step.
func ResumeAt(name string) func(*engine.Spec) {
	return func(spec *engine.Spec) {
		if name == "" {
			return
		}
		for _, step := range spec.Steps {
			if step.Metadata.Name == name {
				return
			}
			if step.Metadata.Name == "clone" {
				continue
			}
			step.RunPolicy = engine.RunNever
		}
	}
}
