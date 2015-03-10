package builder

import (
	"github.com/drone/drone-cli/common"
	"github.com/samalba/dockerclient"
)

// helper function to encode the build step to
// a json string. Primarily used for plugins, which
// expect a json encoded string in stdin or arg[1].
func toCommand(build *Build, step *common.Step) []string {
	payload := BuildPayload{
		build.Repo,
		build.Commit,
		build.Clone,
		step.Config,
	}
	return []string{payload.Encode()}
}

// helper function that converts a build step to
// a hostConfig for use with the dockerclient
func toHostConfig(step *common.Step) *dockerclient.HostConfig {
	return &dockerclient.HostConfig{
		Privileged:  step.Privileged,
		NetworkMode: step.NetworkMode,
	}
}

// helper function that converts the build step to
// a containerConfig for use with the dockerclient
func toContainerConfig(step *common.Step) *dockerclient.ContainerConfig {
	config := &dockerclient.ContainerConfig{
		Image:      step.Name,
		Env:        step.Environment,
		Cmd:        step.Command,
		Entrypoint: step.Entrypoint,
		WorkingDir: step.WorkingDir,
	}

	if len(step.Volumes) != 0 {
		config.Volumes = map[string]struct{}{}
		for _, path := range step.Volumes {
			config.Volumes[path] = struct{}{}
		}
	}

	return config
}
