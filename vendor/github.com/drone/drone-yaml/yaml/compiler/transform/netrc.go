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
	"fmt"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml/yaml/compiler/internal/rand"
)

const (
	netrcName = ".netrc"
	netrcPath = "/var/run/drone/.netrc"
	netrcMode = 0777
)

const disableNetrcMount = true

// WithNetrc is a helper function that creates a netrc file
// and mounts the file to all container steps.
func WithNetrc(machine, username, password string) func(*engine.Spec) {
	return func(spec *engine.Spec) {
		if username == "" || password == "" {
			return
		}

		// TODO(bradrydzewski) temporarily disable mounting
		// the netrc file due to issues with kubernetes
		// compatibility.
		if disableNetrcMount == false {
			// Currently file mounts don't seem to work in Windows so environment
			// variables are used instead
			// FIXME: https://github.com/drone/drone-yaml/issues/20
			if spec.Platform.OS != "windows" {
				netrc := generateNetrc(machine, username, password)
				spec.Files = append(spec.Files, &engine.File{
					Metadata: engine.Metadata{
						UID:       rand.String(),
						Name:      netrcName,
						Namespace: spec.Metadata.Namespace,
					},
					Data: []byte(netrc),
				})
				for _, step := range spec.Steps {
					step.Files = append(step.Files, &engine.FileMount{
						Name: netrcName,
						Path: netrcPath,
						Mode: netrcMode,
					})
				}
			}
		}

		// TODO(bradrydzewski) these should only be injected
		// if the file is not mounted, if OS == Windows.
		for _, step := range spec.Steps {
			if step.Envs == nil {
				step.Envs = map[string]string{}
			}
			step.Envs["CI_NETRC_MACHINE"] = machine
			step.Envs["CI_NETRC_USERNAME"] = username
			step.Envs["CI_NETRC_PASSWORD"] = password

			step.Envs["DRONE_NETRC_MACHINE"] = machine
			step.Envs["DRONE_NETRC_USERNAME"] = username
			step.Envs["DRONE_NETRC_PASSWORD"] = password
		}
	}
}

func generateNetrc(machine, username, password string) string {
	return fmt.Sprintf("machine %s login %s password %s",
		machine, username, password)
}
