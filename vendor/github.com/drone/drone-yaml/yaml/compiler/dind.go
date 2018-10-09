package compiler

import (
	"github.com/drone/drone-yaml/yaml"
	"github.com/drone/drone-yaml/yaml/compiler/image"
)

// DindFunc is a helper function that returns true
// if a container image (specifically a plugin) is
// a whitelisted dind container and should be executed
// in privileged mode.
func DindFunc(images []string) func(*yaml.Container) bool {
	return func(container *yaml.Container) bool {
		// privileged-by-default containers are only
		// enabled for plugins steps that do not define
		// commands.
		if len(container.Commands) > 0 {
			return false
		}
		// privileged-by-default containers are only
		// enabled for plugins steps that do not define
		// custom environment variables. This restriction
		// MAY be lifted in the future.
		if len(container.Environment) > 0 {
			return false
		}
		// if the container image matches any image
		// in the whitelist, return true.
		for _, img := range images {
			a := img
			b := container.Image
			if image.Match(a, b) {
				return true
			}
		}
		return false
	}
}
