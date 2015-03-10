package builder

import (
	"io"

	"github.com/samalba/dockerclient"
)

// Handler defines an interface that can be implemented by
// objects that should be run during the build process.
// to run as part of a build.
type Handler interface {
	Build(*Result) error
	Cancel()
}

type handler struct {
	id     string
	name   string
	detach bool
	client dockerclient.Client
	host   *dockerclient.HostConfig
	config *dockerclient.ContainerConfig
}

func (h *handler) Build(res *Result) error {
	id, err := h.client.CreateContainer(h.config, h.name)
	if err != nil {
		return err
	}
	h.id = id
	err = h.client.StartContainer(h.id, h.host)
	if err != nil {
		return err
	}
	if h.detach {
		return nil
	}
	logs := &dockerclient.LogOptions{
		Follow:     true,
		Stderr:     true,
		Stdout:     true,
		Timestamps: true,
	}
	rc, err := h.client.ContainerLogs(h.id, logs)
	if err != nil {
		return err
	}
	io.Copy(res, rc)
	info, err := h.client.InspectContainer(h.id)
	if err != nil {
		return err
	}
	res.WriteExitCode(info.State.ExitCode)
	return nil
}

func (h *handler) Cancel() {
	h.client.StopContainer(h.id, 5)
	h.client.KillContainer(h.id, "SIGKILL")
	h.client.RemoveContainer(h.id, true, false)
}

// BatchHandler returns a handler that runs a build
// task in batch mode. It will block until the task
// comples, writing stdout to the result.
//
// If the task fails the exit status code is written
// to the result as well.
func BatchHandler(client dockerclient.Client, c *Container) Handler {
	return &handler{
		client: client,
		host:   c.toHostConfig(),
		config: c.toContainerConfig(),
	}
}

// DetachedHandler returns a handler that runs a build
// task in detached mode. It will start the container async,
// and immediately exit.
func DetachedHandler(client dockerclient.Client, c *Container) Handler {
	return &handler{
		detach: true,
		client: client,
		host:   c.toHostConfig(),
		config: c.toContainerConfig(),
	}
}
