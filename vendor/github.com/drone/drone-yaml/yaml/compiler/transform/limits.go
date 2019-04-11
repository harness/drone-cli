// Copyright the Drone Authors.
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

// WithLimits is a transform function that applies
// resource limits to the container processes.
func WithLimits(memlimit, cpulimit int64) func(*engine.Spec) {
	return func(spec *engine.Spec) {
		// if no limits are defined exit immediately.
		if memlimit == 0 && cpulimit == 0 {
			return
		}
		// otherwise apply the resource limits to every
		// step in the runtime spec.
		for _, step := range spec.Steps {
			if step.Resources == nil {
				step.Resources = &engine.Resources{}
			}
			if step.Resources.Limits == nil {
				step.Resources.Limits = &engine.ResourceObject{}
			}
			step.Resources.Limits.Memory = memlimit
			step.Resources.Limits.CPU = cpulimit
		}
	}
}
