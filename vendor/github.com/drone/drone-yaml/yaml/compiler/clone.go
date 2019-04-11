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
	"strconv"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml/yaml"
	"github.com/drone/drone-yaml/yaml/compiler/internal/rand"
)

// default name of the clone step.
const cloneStepName = "clone"

// helper function returns the preferred clone image
// based on the target architecture.
func cloneImage(src *yaml.Pipeline) string {
	switch {
	case src.Platform.OS == "linux" && src.Platform.Arch == "arm":
		return "drone/git:linux-arm"
	case src.Platform.OS == "linux" && src.Platform.Arch == "arm64":
		return "drone/git:linux-arm64"
	case src.Platform.OS == "windows" && src.Platform.Version == "1809":
		return "drone/git:windows-1809-amd64"
	case src.Platform.OS == "windows" && src.Platform.Version == "1803":
		return "drone/git:windows-1803" // TODO update to correct format
	case src.Platform.OS == "windows" && src.Platform.Version == "1709":
		return "drone/git:windows-1709-amd64"
	case src.Platform.OS == "windows":
		return "drone/git:windows-1809-amd64"
	default:
		return "drone/git"
	}
}

// helper function configures the clone depth parameter,
// specific to the clone plugin.
//
// TODO(bradrydzewski) rename to setupCloneParams
func setupCloneDepth(src *yaml.Pipeline, dst *engine.Step) {
	if depth := src.Clone.Depth; depth > 0 {
		dst.Envs["PLUGIN_DEPTH"] = strconv.Itoa(depth)
	}
	if skipVerify := src.Clone.SkipVerify; skipVerify {
		dst.Envs["GIT_SSL_NO_VERIFY"] = "true"
		dst.Envs["PLUGIN_SKIP_VERIFY"] = "true"
	}
}

// helper function configures the .git-clone credentials
// file. The file is mounted into the container, pointed
// to by XDG_CONFIG_HOME
// see https://git-scm.com/docs/git-credential-store
func setupCloneCredentials(spec *engine.Spec, dst *engine.Step, data []byte) {
	if len(data) == 0 {
		return
	}
	// TODO(bradrydzewski) we may need to update the git
	// clone plugin to configure the git credential store.
	dst.Files = append(dst.Files, &engine.FileMount{
		Name: ".git-credentials",
		Path: "/root/.git-credentials",
	})
	spec.Files = append(spec.Files, &engine.File{
		Metadata: engine.Metadata{
			UID:       rand.String(),
			Namespace: spec.Metadata.Namespace,
			Name:      ".git-credentials",
		},
		Data: data,
	})
}

// helper function creates a default container configuration
// for the clone stage. The clone stage is automatically
// added to each pipeline.
func createClone(src *yaml.Pipeline) *yaml.Container {
	return &yaml.Container{
		Name:  cloneStepName,
		Image: cloneImage(src),
	}
}
