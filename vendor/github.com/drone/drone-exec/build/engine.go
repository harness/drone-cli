package build

import (
	"io"

	"github.com/drone/drone-exec/yaml"
)

// Engine defines the container runtime engine.
type Engine interface {
	ContainerStart(*yaml.Container) (string, error)
	ContainerStop(string) error
	ContainerRemove(string) error
	ContainerWait(string) (*State, error)
	ContainerLogs(string) (io.ReadCloser, error)
}
