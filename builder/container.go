package builder

import (
	"github.com/drone/drone-cli/common"
	"github.com/samalba/dockerclient"
)

// Container represents a Docker Container used
// to execute a build step.
type Container struct {
	ID          string
	Image       string
	Pull        bool
	Detached    bool
	Privileged  bool
	Environment []string
	Entrypoint  []string
	Command     []string
	Volumes     []string
	VolumesFrom []string
	WorkingDir  string
	NetworkMode string
}

// helper function to create a container from a step.
func fromStep(step *common.Step) *Container {
	return &Container{
		Image:       step.Name,
		Pull:        step.Pull,
		Privileged:  step.Privileged,
		Volumes:     step.Volumes,
		WorkingDir:  step.WorkingDir,
		NetworkMode: step.NetworkMode,
		Entrypoint:  step.Entrypoint,
		Environment: step.Environment,
		Command:     step.Command,
	}
}

// helper function to create a container from a build
// step. The build task will invoke a shell script
// at an expected path.
func fromBuild(build *Build, step *common.Step) *Container {
	c := fromStep(step)
	c.Entrypoint = []string{"/bin/bash"}
	c.Command = []string{"/drone/bin/build.sh"}
	return c
}

// helper function to create a container from a setup
// step. This is a special container. It is used to
// bootstrap the environment, create build directories,
// and generate the build script.
//
// see https://github.com/drone-plugins/drone-build
func fromSetup(build *Build, step *common.Step) *Container {
	c := fromStep(step)
	c.Image = "plugins/drone-build"
	c.Entrypoint = []string{"/go/bin/drone-build"}
	c.Command = toCommand(build, step)
	return c
}

// helper function to create a container from any plugin
// step, including notification, deployment and publish steps.
// It is used to create the plugin payload (JSON) and pass
// to the container as arg[1]
func fromPlugin(build *Build, step *common.Step) *Container {
	c := fromStep(step)
	c.Entrypoint = []string{}
	c.Command = toCommand(build, step)
	return c
}

// helper function to create a container from a compose
// step. This creates the container almost verbatim. It only
// adds a --detached flag to the container. This instructure
// the build not to block and wait for this container to
// finish execution.
func fromCompose(build *Build, step *common.Step) *Container {
	c := fromStep(step)
	c.Detached = true
	return c
}

// helper function to encode the container arguments
// in a json string. Primarily used for plugins, which
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

// helper function that converts the container to
// a hostConfig for use with the dockerclient
func (c *Container) toHostConfig() *dockerclient.HostConfig {
	return &dockerclient.HostConfig{
		Privileged:  c.Privileged,
		NetworkMode: c.NetworkMode,
	}
}

// helper function that converts the container to
// a containerConfig for use with the dockerclient
func (c *Container) toContainerConfig() *dockerclient.ContainerConfig {
	config := &dockerclient.ContainerConfig{
		Image:      c.Image,
		Env:        c.Environment,
		Cmd:        c.Command,
		Entrypoint: c.Entrypoint,
		WorkingDir: c.WorkingDir,
	}

	if len(c.Volumes) != 0 {
		config.Volumes = map[string]struct{}{}
		for _, path := range c.Volumes {
			config.Volumes[path] = struct{}{}
		}
	}

	return config
}
