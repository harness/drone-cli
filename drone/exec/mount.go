package exec

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dchest/uniuri"
	"github.com/drone/drone-runtime/engine"
)

// helper function mounts the directory (typicall the
// current working directory) as the workspace in all
// runtime containers. The ensure the working directory
// (your .git repository) is available to your local
// execution.
func mountWorkspace(spec *engine.Spec, pwd string) {
	// mount the working directory as a host-machine
	// volume to expose the working directory inside
	// the containers.
	spec.Docker.Volumes = append(
		spec.Docker.Volumes,
		&engine.Volume{
			Metadata: engine.Metadata{
				UID:  uniuri.New(),
				Name: "_local",
			},
			HostPath: &engine.VolumeHostPath{
				Path: pwd,
			},
		},
	)
	// mount the working directory volume into the
	// workspace of every container.
	for _, container := range spec.Steps {
		container.Volumes = append(
			container.Volumes,
			&engine.VolumeMount{
				Name: "_local",
				// HACK(bradrydzewski) this feels like a
				// hack. It would be nice if we could make
				// this an official transform function in
				// the transform.
				Path: container.Envs["DRONE_WORKSPACE"],
			},
		)
	}
}

// helper funciton normalizes the mount path based on the
// host operating system. Specifically we need this for
// windows environments.
func normalizeMountPath(path string) string {
	switch runtime.GOOS {
	case "windows":
		return normalizeMountPathWin(path)
	default:
		return path
	}
}

// helper funciton normalizes the mount path for windows
// environments.
func normalizeMountPathWin(path string) string {
	base := filepath.VolumeName(path)
	if len(base) == 2 {
		path = strings.TrimPrefix(path, base)
		base = strings.ToLower(strings.TrimSuffix(base, ":"))
		return "/" + base + filepath.ToSlash(path)
	}
	return filepath.ToSlash(path)
}
