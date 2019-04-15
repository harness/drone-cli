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

// WithAuths is a transform function that adds a set
// of global registry credentials to the container.
func WithAuths(auths []*engine.DockerAuth) func(*engine.Spec) {
	return func(spec *engine.Spec) {
		spec.Docker.Auths = append(spec.Docker.Auths, auths...)
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
