package docker

import (
	"io"

	"github.com/drone/drone-cli/engine"
	"github.com/samalba/dockerclient"
)

type Docker struct {
	client     dockerclient.Client
	ambassador *engine.Container
}

func New(client dockerclient.Client) *Docker {
	return &Docker{client: client}
}

func (d *Docker) Create(c *engine.Container) error {
	return create(d.client, c)
}

func (d *Docker) Start(c *engine.Container) error {
	return start(d.client, d.ambassador, c)
}

func (d *Docker) Stop(c *engine.Container) error {
	return stop(d.client, c)
}

func (d *Docker) Remove(c *engine.Container) error {
	return remove(d.client, c)
}

func (d *Docker) State(c *engine.Container) (*engine.State, error) {
	return state(d.client, c)
}

func (d *Docker) Logs(c *engine.Container) (io.ReadCloser, error) {
	return logs(d.client, c)
}

func (d *Docker) Wait(c *engine.Container) error {
	return wait(d.client, c)
}

func (d *Docker) Setup() error {
	d.ambassador = &engine.Container{
		Image:      "busybox",
		Volumes:    []string{"/drone"},
		Entrypoint: []string{"/bin/sleep", "1d"},
	}
	err := create(d.client, d.ambassador)
	if err != nil {
		return err
	}
	return start(d.client, nil, d.ambassador)
}

func (d *Docker) Teardown() error {
	stop(d.client, d.ambassador)
	return remove(d.client, d.ambassador)
}
