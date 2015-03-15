package common

// Container is a running instance
type Container struct {
	ID    string
	Type  string
	Image *Image

	// Engine is the engine that is runnnig the container
	//Engine *Engine
}
