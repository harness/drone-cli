package transform

import (
	"fmt"

	"github.com/drone/drone-runtime/engine"
)

const (
	netrcName = ".netrc"
	netrcPath = "/root/.netrc"
	netrcMode = 0600
)

// WithNetrc is a helper function that creates a netrc file
// and mounts the file to all container steps.
func WithNetrc(machine, username, password string) func(*engine.Spec) {
	return func(spec *engine.Spec) {
		if username == "" || password == "" {
			return
		}
		netrc := generateNetrc(machine, username, password)
		spec.Files = append(spec.Files, &engine.File{
			Name: netrcName,
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

func generateNetrc(machine, username, password string) string {
	return fmt.Sprintf("machine %s login %s password %s",
		machine, username, password)
}
