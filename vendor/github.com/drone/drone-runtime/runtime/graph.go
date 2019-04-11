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

package runtime

import "github.com/drone/drone-runtime/engine"

type status struct {
	step  *engine.Step
	state *engine.State
}

// isSerial returns true if the steps are to be executed
// in serial mode, with no graph dependencies defined.
func isSerial(spec *engine.Spec) bool {
	for _, step := range spec.Steps {
		// if a single dependency is defined we can exit
		// and return false.
		if len(step.DependsOn) != 0 {
			return false
		}
	}
	return true
}

// nextStep returns the next step in the dependency graph.
// If no steps are ready for execution, a nil value is
// returned.
func nextStep(spec *engine.Spec, states map[string]*status) *engine.Step {
LOOP:
	for _, step := range spec.Steps {
		// if the step has already stated, move to the
		// next step in the list.
		state := states[step.Metadata.Name]
		if state.state != nil {
			continue
		}

		// if the step has zero dependencies and has not
		// started, it can be started immediately.
		if len(step.DependsOn) == 0 {
			return step
		}
		// if the step has dependencies, we check to ensure
		// all dependent steps are complete. If no, move on
		// to test the next step.
		for _, name := range step.DependsOn {
			state, ok := states[name]
			// if the dependency does not exist in the
			// state map it is considered fulfilled to
			// avoid deadlock.
			if !ok {
				continue
			}

			// if the dependency is running in detached
			// mode it is considered fulfilled to avoid
			// deadlock.
			if state.step.Detach {
				continue
			}
			// if the dependency is skipped (never executed)
			// it is considered fulfilled to avoid deadlock.
			if state.step.RunPolicy == engine.RunNever {
				continue
			}
			// if the dependency has not executed, the step
			// is not ready for execution. Break to the
			// next step in the specification list.
			if state.state == nil || state.state.Exited == false {
				break LOOP
			}
		}
		// if all dependencies are completed, the step
		// can be started.
		return step
	}
	return nil
}
