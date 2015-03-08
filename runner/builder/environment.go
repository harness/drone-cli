package builder

import (
	"github.com/samalba/dockerclient"
)

// Default image used when creating the ambassador.
const ambassadorImage = "busybox"

// Default disk volume when creating the ambassador,
// shared across the group of containers
var ambassadorVolume = []string{"/drone"}

// Default command used to start the amabassador contianer
// and block until killed.
var ambassadorEntry = []string{"/bin/sleep", "1d"}

// An Environment defines the disk and network environment
// through and embassador container.
type Environment struct {
	client  dockerclient.Client
	command *Command
}

// NewEnvironment creates a new environment
func NewEnvironment(client dockerclient.Client) *Environment {
	return &Environment{client: client}
}

// Setup starts the ambassador container.
func (e *Environment) Setup() error {
	command := &Command{
		Image:      ambassadorImage,
		Entrypoint: ambassadorEntry,
		Volumes:    ambassadorVolume,
	}
	err := create(command, e.client)
	if err != nil {
		return err
	}
	config := &dockerclient.HostConfig{}
	return e.client.StartContainer(command.ID, config)
}

// Destroy destroys the ambassador container.
func (e *Environment) Destroy() error {
	e.client.StopContainer(e.command.ID, 5)
	e.client.KillContainer(e.command.ID, "SIGKILL")
	return e.client.RemoveContainer(e.command.ID, true, true)
}

func (e *Environment) String() string {
	return "container:" + e.command.ID
}
