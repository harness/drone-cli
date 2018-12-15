package compiler

import (
	unixpath "path"
	"strings"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml/yaml"
)

// default name for the workspace volume.
const workspaceName = "workspace"

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

// helper function mounts the working directory base
// path to the container.
func setupWorkingDirMount(step *engine.Step, path string) {
	step.Volumes = append(step.Volumes, &engine.VolumeMount{
		Name: workspaceName,
		Path: path,
	})
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
		base = "/drone"
	}
	if path == "" {
		path = "src"
	}
	full = unixpath.Join(base, path)

	if from.Platform.OS == "windows" {
		base = toWindowsDrive(base)
		path = toWindowsPath(path)
		full = toWindowsDrive(full)
	}
	return base, path, full
}
