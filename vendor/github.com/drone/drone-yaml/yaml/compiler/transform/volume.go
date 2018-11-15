package transform

import (
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
