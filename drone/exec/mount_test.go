package exec

import (
	"testing"

	"github.com/drone/drone-runtime/engine"
)

func TestMountWorkspace(t *testing.T) {
	step := &engine.Step{
		Envs: map[string]string{
			"DRONE_WORKSPACE": "/workspace",
		},
	}
	spec := &engine.Spec{
		Docker: &engine.DockerConfig{},
		Steps:  []*engine.Step{step},
	}
	mountWorkspace(spec, "/path/on/host")

	if len(spec.Docker.Volumes) == 0 {
		t.Errorf("Expect volume mounted on host")
		return
	}

	volume := spec.Docker.Volumes[0]
	if got, want := volume.HostPath.Path, "/path/on/host"; got != want {
		t.Errorf("Want volume mount %s, got %s", want, got)
	}
	if got, want := volume.Metadata.Name, "_local"; got != want {
		t.Errorf("Want volume name %s, got %s", want, got)
	}
	if volume.Metadata.UID == "" {
		t.Errorf("Want volume UID got empty string")
	}

	if len(step.Volumes) == 0 {
		t.Errorf("Excpect container volume mount")
		return
	}
	mount := step.Volumes[0]
	if got, want := mount.Path, "/workspace"; got != want {
		t.Errorf("Want volume mount %s, got %s", want, got)
	}
	if got, want := mount.Name, "_local"; got != want {
		t.Errorf("Want volume name %s, got %s", want, got)
	}
}
