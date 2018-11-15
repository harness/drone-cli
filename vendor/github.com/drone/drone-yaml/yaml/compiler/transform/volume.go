package transform

import (
	"github.com/dchest/uniuri"
	"github.com/drone/drone-runtime/engine"
)

// WithVolumes is a transform function that adds a set
// of global volumes to the container.
func WithVolumes(volumes map[string]string) func(*engine.Spec) {
	return func(spec *engine.Spec) {
		for key, value := range volumes {
			volume := &engine.Volume{
				Metadata: engine.Metadata{
					UID:       uniuri.New(),
					Name:      uniuri.New(),
					Namespace: spec.Metadata.Namespace,
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
