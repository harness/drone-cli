package runner

import (
	"io"
	"io/ioutil"

	"github.com/samalba/dockerclient"
)

const (
	ContainerTypeBatch      = "batch"
	ContainerTypeService    = "service"
	ContainerTypeAmbassador = "ambassador"
)

type Container struct {
	Name        string
	Image       string
	Pull        bool
	Detached    bool
	Privileged  bool
	Env         []string
	Cmd         []string
	Entrypoint  []string
	WorkingDir  string
	NetworkMode string
	Volumes     []string
	VolumesFrom []string
	Links       []string

	client dockerclient.Client
	info   *dockerclient.ContainerInfo
}

func (c *Container) SetClient(client dockerclient.Client) {
	c.client = client
}

func (c *Container) Create() error {
	config := dockerclient.ContainerConfig{
		Image:      c.Image,
		Env:        c.Env,
		Cmd:        c.Cmd,
		Entrypoint: c.Entrypoint,
		WorkingDir: c.WorkingDir,
	}

	if len(c.Volumes) != 0 {
		config.Volumes = map[string]struct{}{}
		for _, path := range c.Volumes {
			config.Volumes[path] = struct{}{}
		}
	}

	// TODO: we should inspect the image first. We should
	// run `docker pull` if inspect returns an error or pull == true
	if c.Pull {
		if err := c.client.PullImage(c.Image, nil); err != nil {
			return err
		}
	}

	id, err := c.client.CreateContainer(&config, c.Name)
	if err != nil {
		return err
	}
	c.info, err = c.client.InspectContainer(id)
	return err
}

func (c *Container) Start() error {
	config := dockerclient.HostConfig{
		Privileged:  c.Privileged,
		NetworkMode: c.NetworkMode,
		VolumesFrom: c.VolumesFrom,
	}
	return c.client.StartContainer(c.info.Id, &config)
}

func (c *Container) Stop() error {
	return c.client.StopContainer(c.info.Id, 10)
}

func (c *Container) Kill() error {
	return c.client.KillContainer(c.info.Id, "SIGKILL")
}

func (c *Container) Remove() error {
	return c.client.RemoveContainer(c.info.Id, true, false)
}

func (c *Container) RemoveVolume() error {
	return c.client.RemoveContainer(c.info.Id, true, true)
}

func (c *Container) Wait() error {
	src, err := c.Logs()
	if err != nil {
		return err
	}
	defer src.Close()
	_, err = io.Copy(ioutil.Discard, src)
	return nil
}

func (c *Container) Logs() (io.ReadCloser, error) {
	opts := dockerclient.LogOptions{
		Follow:     true,
		Stderr:     true,
		Stdout:     true,
		Tail:       10000,
		Timestamps: true,
	}
	return c.client.ContainerLogs(c.info.Id, &opts)
}

func (c *Container) Inspect() (*dockerclient.ContainerInfo, error) {
	return c.client.InspectContainer(c.info.Id)
}
