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

package compiler

import (
	"github.com/drone/drone-yaml/yaml"
	"github.com/drone/drone-yaml/yaml/compiler/image"
)

// DindFunc is a helper function that returns true
// if a container image (specifically a plugin) is
// a whitelisted dind container and should be executed
// in privileged mode.
func DindFunc(images []string) func(*yaml.Container) bool {
	return func(container *yaml.Container) bool {
		// privileged-by-default containers are only
		// enabled for plugins steps that do not define
		// commands, command, or entrypoint.
		if len(container.Commands) > 0 {
			return false
		}
		if len(container.Command) > 0 {
			return false
		}
		if len(container.Entrypoint) > 0 {
			return false
		}
		// if the container image matches any image
		// in the whitelist, return true.
		for _, img := range images {
			a := img
			b := container.Image
			if image.Match(a, b) {
				return true
			}
		}
		return false
	}
}
