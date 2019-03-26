// Copyright 2019 Drone IO, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
