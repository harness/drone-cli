package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/samalba/dockerclient"
)

var errNop = errors.New("Operation not supported")

type mockClient struct{}

func newMockClient() (dockerclient.Client, error) {
	return &mockClient{}, nil
}

// CreateContainer creates a container.
func (c *mockClient) CreateContainer(config *dockerclient.ContainerConfig, name string) (string, error) {
	fmt.Println("CREATE --------------------------------------------------------")
	fmt.Println("  Image :", config.Image)
	fmt.Println("  Env   :", config.Env)
	fmt.Println("  Cmd   :", config.Cmd)
	fmt.Println("  Entry :", config.Entrypoint)
	return config.Image, nil
}

// InspectContainer returns container details.
func (c *mockClient) InspectContainer(id string) (*dockerclient.ContainerInfo, error) {
	info := dockerclient.ContainerInfo{}
	return &info, nil
}

// ContainerLogs returns an io.ReadCloser for reading the
// container logs.
func (c *mockClient) ContainerLogs(id string, options *dockerclient.LogOptions) (io.ReadCloser, error) {
	var buf bytes.Buffer
	buf.WriteString("")
	return ioutil.NopCloser(&buf), nil
}

// StartContainer starts a container. The mockClient volume
// is automatically linked. The mockClient network is linked
// iff a network mode is not already specified.
func (c *mockClient) StartContainer(id string, config *dockerclient.HostConfig) error {
	fmt.Println("START ---------------------------------------------------------")
	fmt.Println("  ID   :", id)
	fmt.Println("  Net  :", config.NetworkMode)
	fmt.Println("  Priv :", config.Privileged)
	fmt.Println("  Vols :", config.VolumesFrom)
	return nil
}

// StopContainer stops a container.
func (c *mockClient) StopContainer(id string, timeout int) error {
	fmt.Println("STOP ----------------------------------------------------------")
	fmt.Println("  ID :", id)
	return nil
}

// PullImage pulls an image.
func (c *mockClient) PullImage(name string, auth *dockerclient.AuthConfig) error {
	return nil
}

// RemoveContainer removes a container.
func (c *mockClient) RemoveContainer(id string, force, volumes bool) error {
	fmt.Println("REMOVE --------------------------------------------------------")
	fmt.Println("  ID :", id)
	return nil
}

// KillContainer kills a running container.
func (c *mockClient) KillContainer(id, signal string) error {
	fmt.Println("KILL ----------------------------------------------------------")
	fmt.Println("  ID  :", id)
	fmt.Println("  SIG :", signal)
	return nil
}

//
// methods below are not implemented
//

// Info returns a no-op error
func (c *mockClient) Info() (*dockerclient.Info, error) {
	return nil, errNop
}

// ListContainers returns a no-op error
func (c *mockClient) ListContainers(all bool, size bool, filters string) ([]dockerclient.Container, error) {
	return nil, errNop
}

// RestartContainer returns a no-op error
func (c *mockClient) RestartContainer(id string, timeout int) error {
	return errNop
}

// StartMonitorEvents returns a no-op error
func (c *mockClient) StartMonitorEvents(cb dockerclient.Callback, ec chan error, args ...interface{}) {

}

// StopAllMonitorEvents returns a no-op error
func (c *mockClient) StopAllMonitorEvents() {

}

// Version returns a no-op error
func (c *mockClient) Version() (*dockerclient.Version, error) {
	return nil, errNop
}

// ListImages returns a no-op error
func (c *mockClient) ListImages() ([]*dockerclient.Image, error) {
	return nil, errNop
}

// RemoveImage returns a no-op error
func (c *mockClient) RemoveImage(name string) error {
	return errNop
}

// PauseContainer returns a no-op error
func (c *mockClient) PauseContainer(name string) error {
	return errNop
}

// UnpauseContainer returns a no-op error
func (c *mockClient) UnpauseContainer(name string) error {
	return errNop
}

// Exec returns a no-op error
func (c *mockClient) Exec(config *dockerclient.ExecConfig) (string, error) {
	var empty string
	return empty, errNop
}
