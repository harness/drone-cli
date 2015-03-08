package cluster

import "io"

type Cluster interface {
	Create(*Container) error
	Remove(*Container) error
	Start(*Container) error
	Stop(*Container) error
	Wait(*Container) error
	Logs(*Container) (io.ReadCloser, error)
	State(*Container) (*State, error)
}
