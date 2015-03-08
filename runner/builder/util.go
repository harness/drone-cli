package builder

import (
	"io"
	"io/ioutil"

	"github.com/samalba/dockerclient"
)

// helper function to create a container from a
// build command.
func create(c *Command, client dockerclient.Client) error {
	config := dockerclient.ContainerConfig{
		Image:      c.Image,
		Cmd:        c.Cmd,
		Env:        c.Env,
		Entrypoint: c.Entrypoint,
	}

	if len(c.Volumes) != 0 {
		config.Volumes = map[string]struct{}{}
		for _, path := range c.Volumes {
			config.Volumes[path] = struct{}{}
		}
	}

	var err error
	c.ID, err = client.CreateContainer(&config, "")
	return err
}

// helper function to wait until a container stops.
// TODO (brydzewski) send patch for dockerclient wait function
func wait(c *Command, client dockerclient.Client) error {
	src, err := logs(c, client)
	if err != nil {
		return err
	}
	defer src.Close()
	_, err = io.Copy(ioutil.Discard, src)
	return err
}

// helper function to stream a container logs.
func logs(c *Command, client dockerclient.Client) (io.ReadCloser, error) {
	return client.ContainerLogs(c.ID, &dockerclient.LogOptions{
		Follow:     true,
		Stderr:     true,
		Stdout:     true,
		Timestamps: true,
	})
}
