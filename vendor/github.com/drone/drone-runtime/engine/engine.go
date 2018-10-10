package engine

//go:generate mockgen -source=engine.go -destination=mocks/engine.go

import (
	"context"
	"io"
)

// Engine defines a runtime engine for pipeline execution.
type Engine interface {
	// Setup the pipeline environment.
	Setup(context.Context, *Spec) error

	// Create creates the pipeline state.
	Create(context.Context, *Spec, *Step) error

	// Start the pipeline step.
	Start(context.Context, *Spec, *Step) error

	// Wait for the pipeline step to complete and returns the completion results.
	Wait(context.Context, *Spec, *Step) (*State, error)

	// Tail the pipeline step logs.
	Tail(context.Context, *Spec, *Step) (io.ReadCloser, error)

	// Destroy the pipeline environment.
	Destroy(context.Context, *Spec) error
}
