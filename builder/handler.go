package builder

import (
	"io"

	"github.com/drone/drone-cli/common"
	"github.com/samalba/dockerclient"
)

// Handler defines an interface that can be implemented by objects
// that should be run during the build process. to run as part of a build.
type Handler interface {
	Build(*ResultWriter) error
	Cancel()
}

type handler struct {
	name   string
	pull   bool
	detach bool
	client dockerclient.Client
	host   *dockerclient.HostConfig
	config *dockerclient.ContainerConfig
}

func (h *handler) Build(rw *ResultWriter) error {
	name, err := h.client.CreateContainer(h.config, "")
	if err != nil {
		// on error try to pull the Docker image.
		// note that this may not be the cause of
		// the error, but we'll try just in case.
		h.client.PullImage(h.config.Image, nil)

		// then try to re-create
		name, err = h.client.CreateContainer(h.config, "")
		if err != nil {
			return err
		}
	}
	h.name = name
	err = h.client.StartContainer(h.name, h.host)
	if err != nil {
		return err
	}

	// if we are running a service container we
	// can just exit right here.
	if h.detach {
		return nil
	}

	// get the docker logs and write to the resposne.
	logs := &dockerclient.LogOptions{
		Follow:     true,
		Stderr:     true,
		Stdout:     true,
		Timestamps: true,
	}
	rc, err := h.client.ContainerLogs(h.name, logs)
	if err != nil {
		return err
	}
	io.Copy(rw, rc)

	// get the container state and write the exit status
	// to the response.
	info, err := h.client.InspectContainer(h.name)
	if err != nil {
		return err
	}
	rw.WriteExitCode(info.State.ExitCode)
	return nil
}

func (h *handler) Cancel() {
	h.client.StopContainer(h.name, 5)
	h.client.KillContainer(h.name, "SIGKILL")
	h.client.RemoveContainer(h.name, true, false)
}

// Batch returns a handler that launches a container. The container
// will start and block until the container exits.
//
// The container output and result are written to the ResponseWriter.
func Batch(build *Build, step *common.Step) Handler {
	host := toHostConfig(step)
	conf := toContainerConfig(step)
	if step.Config != nil {
		conf.Cmd = toCommand(build, step)
		conf.Entrypoint = []string{}
	}
	return &handler{
		pull:   step.Pull,
		client: build.Client,
		config: conf,
		host:   host,
	}
}

// Script returns a handler that launches the build script
// container. The setup or bootstrap container is a pre-requisite.
//
// The shell script generated in the setup or bootstrap step
// will be used as the container entrypoint.
func Script(build *Build, step *common.Step) Handler {
	host := toHostConfig(step)
	conf := toContainerConfig(step)
	conf.Entrypoint = []string{"/bin/bash"}
	conf.Cmd = []string{"/drone/bin/build.sh"}
	return &handler{
		client: build.Client,
		config: conf,
		host:   host,
	}
}

// Setup returns a handler that launches a special type of
// container used to setup or bootstrap a build environment.
//
// This container will setup the project workspace and generate
// the build script.
func Setup(build *Build, step *common.Step) Handler {
	setup := &common.Step{
		Image: "plugins/drone-build",
	}
	host := toHostConfig(setup)
	conf := toContainerConfig(setup)
	conf.Cmd = toCommand(build, step)
	return &handler{
		client: build.Client,
		config: conf,
		host:   host,
	}
}

// Service returns a handler that launches a service container as
// a daemon. The the container will start and then exit immediately,
// without blocking.
//
// The output and exit status will not be written to the ResponseWriter.
func Service(build *Build, step *common.Step) Handler {
	host := toHostConfig(step)
	conf := toContainerConfig(step)
	return &handler{
		pull:   step.Pull,
		client: build.Client,
		detach: true,
		config: conf,
		host:   host,
	}
}
