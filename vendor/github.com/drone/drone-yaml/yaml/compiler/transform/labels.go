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
