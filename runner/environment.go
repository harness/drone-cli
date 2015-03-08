package runner

import (
	"github.com/samalba/dockerclient"
)

const (
	ambassadorImage  = "busybox"
	ambassadorVolume = "/drone"
)

// Default command used to start the amabassador contianer
// and block until killed.
var ambassadorEntry = []string{"/bin/sleep", "1d"}

// An Environment represents a build exeuction environment used
// to execute a single build.
type Environment struct {
	client     dockerclient.Client
	ambassador *dockerclient.ContainerInfo
	containers []*dockerclient.ContainerInfo
}

// NewEnvironment creates a new environment
func NewEnvironment(client dockerclient.Client) *Environment {
	return &Environment{client: client}
}

// Prepare prepares the build execution environment.
func (e *Environment) Prepare() error {
	config := dockerclient.ContainerConfig{
		Image:      ambassadorImage,
		Entrypoint: ambassadorEntry,
	}
	config.Volumes = map[string]struct{}{}
	config.Volumes[ambassadorVolume] = struct{}{}
	id, err := e.client.CreateContainer(&config, "")
	if err != nil {
		return err
	}
	e.ambassador, err = e.client.InspectContainer(id)
	if err != nil {
		return err
	}
	return e.client.StartContainer(e.ambassador.Id, nil)
}

func (e *Environment) Start() error {
	return nil
}

// Destroy destroys the build execution environment.
func (e *Environment) Destroy() error {
	for _, container := range e.containers {
		e.client.StopContainer(container.Id, 5)
		e.client.KillContainer(container.Id, "SIGKILL")
		e.client.RemoveContainer(container.Id, true, false)
	}
	e.client.StopContainer(e.ambassador.Id, 5)
	e.client.KillContainer(e.ambassador.Id, "SIGKILL")
	return e.client.RemoveContainer(e.ambassador.Id, true, true)
}
