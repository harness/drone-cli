package runtime

import "github.com/drone/drone-runtime/engine"

// State defines the pipeline and process state.
type State struct {
	hook   *Hook
	config *engine.Spec
	engine engine.Engine

	// Global state of the runtime.
	Runtime struct {
		// Runtime time started
		Time int64

		// Runtime pipeline error state
		Error error
	}

	// Runtime pipeline step
	Step *engine.Step

	// Current process state.
	State *engine.State
}

// snapshot makes a snapshot of the runtime state.
func snapshot(r *Runtime, step *engine.Step, state *engine.State) *State {
	s := new(State)
	s.Runtime.Error = r.error
	s.Runtime.Time = r.start
	s.config = r.config
	s.hook = r.hook
	s.engine = r.engine
	s.Step = step
	s.State = state
	return s
}
