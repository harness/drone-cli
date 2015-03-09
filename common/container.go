package common

// Container is a running instance
type Container struct {
	// ID is the container's id
	ID string

	// Image is the configuration from which the container was created
	Image *Image

	// Engine is the engine that is runnnig the container
	//Engine *Engine
}
