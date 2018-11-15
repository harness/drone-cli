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
	case src.Platform.OS == "windows":
		return "drone/git:windows-1803"
	default:
		return "drone/git"
	}
}

// helper function configures the clone depth parameter,
// specific to the clone plugin.
func setupCloneDepth(src *yaml.Pipeline, dst *engine.Step) {
	if depth := src.Clone.Depth; depth > 0 {
		dst.Envs["PLUGIN_DEPTH"] = strconv.Itoa(depth)
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
