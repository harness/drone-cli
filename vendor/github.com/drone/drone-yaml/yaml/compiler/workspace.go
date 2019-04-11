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
	unixpath "path"
	"strings"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml/yaml"
	"github.com/drone/drone-yaml/yaml/compiler/internal/rand"
)

const (
	workspacePath     = "/drone/src"
	workspaceName     = "workspace"
	workspaceHostName = "host"
)

func setupWorkingDir(src *yaml.Container, dst *engine.Step, path string) {
	// if the working directory is already set
	// do not alter.
	if dst.WorkingDir != "" {
		return
	}
	// if the user is running the container as a
	// service (detached mode) with no commands, we
	// should use the default working directory.
	if dst.Detach && len(src.Commands) == 0 {
		return
	}
	// else set the working directory.
	dst.WorkingDir = path
}

// helper function appends the workspace base and
// path to the step's list of environment variables.
func setupWorkspaceEnv(step *engine.Step, base, path, full string) {
	step.Envs["DRONE_WORKSPACE_BASE"] = base
	step.Envs["DRONE_WORKSPACE_PATH"] = path
	step.Envs["DRONE_WORKSPACE"] = full
	step.Envs["CI_WORKSPACE_BASE"] = base
	step.Envs["CI_WORKSPACE_PATH"] = path
	step.Envs["CI_WORKSPACE"] = full
}

// helper function converts the path to a valid windows
// path, including the default C drive.
func toWindowsDrive(s string) string {
	return "c:" + toWindowsPath(s)
}

// helper function converts the path to a valid windows
// path, replacing backslashes with forward slashes.
func toWindowsPath(s string) string {
	return strings.Replace(s, "/", "\\", -1)
}

//
//
//

func createWorkspace(from *yaml.Pipeline) (base, path, full string) {
	base = from.Workspace.Base
	path = from.Workspace.Path
	if base == "" {
		base = workspacePath
	}
	full = unixpath.Join(base, path)

	if from.Platform.OS == "windows" {
		base = toWindowsDrive(base)
		path = toWindowsPath(path)
		full = toWindowsDrive(full)
	}
	return base, path, full
}

//
//
//

// CreateWorkspace creates the workspace volume as
// an empty directory mount.
func CreateWorkspace(spec *engine.Spec) {
	spec.Docker.Volumes = append(spec.Docker.Volumes,
		&engine.Volume{
			Metadata: engine.Metadata{
				UID:       rand.String(),
				Name:      workspaceName,
				Namespace: spec.Metadata.Namespace,
				Labels:    map[string]string{},
			},
			EmptyDir: &engine.VolumeEmptyDir{},
		},
	)
}

// CreateHostWorkspace returns a WorkspaceFunc that
// mounts a host machine volume as the pipeline
// workspace.
func CreateHostWorkspace(workdir string) func(*engine.Spec) {
	return func(spec *engine.Spec) {
		CreateWorkspace(spec)
		spec.Docker.Volumes = append(
			spec.Docker.Volumes,
			&engine.Volume{
				Metadata: engine.Metadata{
					UID:  rand.String(),
					Name: workspaceHostName,
				},
				HostPath: &engine.VolumeHostPath{
					Path: workdir,
				},
			},
		)
	}
}

//
//
//

// MountWorkspace is a WorkspaceFunc that mounts the
// default workspace volume to the pipeline step.
func MountWorkspace(step *engine.Step, base, path, full string) {
	step.Volumes = append(step.Volumes, &engine.VolumeMount{
		Name: workspaceName,
		Path: base,
	})
}

// MountHostWorkspace is a WorkspaceFunc that mounts
// the default workspace and host volume to the pipeline.
func MountHostWorkspace(step *engine.Step, base, path, full string) {
	step.Volumes = append(step.Volumes, &engine.VolumeMount{
		Name: workspaceHostName,
		Path: full,
	})
	if path != "" {
		step.Volumes = append(step.Volumes, &engine.VolumeMount{
			Name: workspaceName,
			Path: base,
		})
	}
}
