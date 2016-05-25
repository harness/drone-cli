package client

import (
	"io"

	"github.com/drone/drone-go/drone"
)

// Client is used to communicate with a Drone server.
type Client interface {

	// Pull pulls work from the server queue.
	Pull(os, arch string) (*drone.Payload, error)

	// Push pushes an update to the server.
	Push(*drone.Payload) error

	// Stream streams the build logs to the server.
	Stream(int64, io.ReadCloser) error

	LogStream(int64) (StreamWriter, error)

	LogPost(int64, io.ReadCloser) error

	// Wait waits for the job to the complete.
	Wait(int64) *Wait

	// Ping the server
	Ping() error
}
