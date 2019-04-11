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

import (
	"strings"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml/yaml/compiler/internal/rand"
)

// WithVolumes is a transform function that adds a set
// of global volumes to the container.
func WithVolumes(volumes map[string]string) func(*engine.Spec) {
	return func(spec *engine.Spec) {
		for key, value := range volumes {
			volume := &engine.Volume{
				Metadata: engine.Metadata{
					UID:       rand.String(),
					Name:      rand.String(),
					Namespace: spec.Metadata.Name,
					Labels:    map[string]string{},
				},
				HostPath: &engine.VolumeHostPath{
					Path: key,
				},
			}
			spec.Docker.Volumes = append(spec.Docker.Volumes, volume)
			for _, step := range spec.Steps {
				mount := &engine.VolumeMount{
					Name: volume.Metadata.Name,
					Path: value,
				}
				step.Volumes = append(step.Volumes, mount)
			}
		}
	}
}

// WithVolumeSlice is a transform function that adds a set
// of global volumes to the container that are defined in
// --volume=host:container format.
func WithVolumeSlice(volumes []string) func(*engine.Spec) {
	to := map[string]string{}
	for _, s := range volumes {
		parts := strings.Split(s, ":")
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		val := parts[1]
		to[key] = val
	}
	return WithVolumes(to)
}
