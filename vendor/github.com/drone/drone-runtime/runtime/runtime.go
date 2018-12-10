package runtime

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/drone/drone-runtime/engine"
	"github.com/natessilva/dag"
	"golang.org/x/sync/errgroup"
)

// Runtime executes a pipeline configuration.
type Runtime struct {
	mu sync.Mutex

	engine engine.Engine
	config *engine.Spec
	hook   *Hook
	start  int64
	error  error
}

// New returns a new runtime using the specified runtime
// configuration and runtime engine.
func New(opts ...Option) *Runtime {
	r := &Runtime{}
	r.hook = &Hook{}
	for _, opts := range opts {
		opts(r)
	}
	return r
}

// Run starts the pipeline and waits for it to complete.
func (r *Runtime) Run(ctx context.Context) error {
	return r.Resume(ctx, 0)
}

// Resume starts the pipeline at the specified stage and
// waits for it to complete.
func (r *Runtime) Resume(ctx context.Context, start int) error {
	defer func() {
		// note that we use a new context to destroy the
		// environment to ensure it is not in a canceled
		// state.
		r.engine.Destroy(context.Background(), r.config)
	}()

	r.error = nil
	r.start = time.Now().Unix()

	if r.hook.Before != nil {
		state := snapshot(r, nil, nil)
		if err := r.hook.Before(state); err != nil {
			return err
		}
	}

	if err := r.engine.Setup(ctx, r.config); err != nil {
		return err
	}

	if isSerial(r.config) {
		for i, step := range r.config.Steps {
			steps := []*engine.Step{step}
			if i < start {
				continue
			}
			select {
			case <-ctx.Done():
				return ErrCancel
			case err := <-r.execAll(steps):
				if err != nil {
					r.error = err
				}
			}
		}
	} else {
		err := r.execGraph(ctx)
		if err != nil {
			return err
		}
	}

	if r.hook.After != nil {
		state := snapshot(r, nil, nil)
		if err := r.hook.After(state); err != nil {
			return err
		}
	}
	return r.error
}

func (r *Runtime) execGraph(ctx context.Context) error {
	var d dag.Runner
	for _, s := range r.config.Steps {
		step := s
		d.AddVertex(step.Metadata.Name, func() error {
			select {
			case <-ctx.Done():
				return ErrCancel
			default:
			}
			err := r.exec(step)
			if err != nil {
				r.mu.Lock()
				r.error = err
				r.mu.Unlock()
			}
			return nil
		})
	}
	for _, s := range r.config.Steps {
		for _, dep := range s.DependsOn {
			d.AddEdge(dep, s.Metadata.Name)
		}
	}
	return d.Run()
}

func (r *Runtime) execAll(group []*engine.Step) <-chan error {
	var g errgroup.Group
	done := make(chan error)

	for _, step := range group {
		step := step
		g.Go(func() error {
			return r.exec(step)
		})
	}

	go func() {
		done <- g.Wait()
		close(done)
	}()
	return done
}

func (r *Runtime) exec(step *engine.Step) error {
	ctx := context.Background()

	switch {
	case step.RunPolicy == engine.RunNever:
		return nil
	case r.error != nil && step.RunPolicy == engine.RunOnSuccess:
		return nil
	case r.error == nil && step.RunPolicy == engine.RunOnFailure:
		return nil
	}

	if r.hook.BeforeEach != nil {
		state := snapshot(r, step, nil)
		if err := r.hook.BeforeEach(state); err == ErrSkip {
			return nil
		} else if err != nil {
			return err
		}
	}

	if err := r.engine.Create(ctx, r.config, step); err != nil {
		// TODO(bradrydzewski) refactor duplicate code
		if r.hook.AfterEach != nil {
			r.hook.AfterEach(
				snapshot(r, step, &engine.State{
					ExitCode: 255, Exited: true,
				}),
			)
		}
		return err
	}

	if err := r.engine.Start(ctx, r.config, step); err != nil {
		// TODO(bradrydzewski) refactor duplicate code
		if r.hook.AfterEach != nil {
			r.hook.AfterEach(
				snapshot(r, step, &engine.State{
					ExitCode: 255, Exited: true,
				}),
			)
		}
		return err
	}

	rc, err := r.engine.Tail(ctx, r.config, step)
	if err != nil {
		// TODO(bradrydzewski) refactor duplicate code
		if r.hook.AfterEach != nil {
			r.hook.AfterEach(
				snapshot(r, step, &engine.State{
					ExitCode: 255, Exited: true,
				}),
			)
		}
		return err
	}

	var g errgroup.Group
	state := snapshot(r, step, nil)
	g.Go(func() error {
		return stream(state, rc)
	})

	if step.Detach {
		return nil // do not wait for service containers to complete.
	}

	defer func() {
		g.Wait() // wait for background tasks to complete.
		rc.Close()
	}()

	wait, err := r.engine.Wait(ctx, r.config, step)
	if err != nil {
		return err
	}

	err = g.Wait() // wait for background tasks to complete.

	if wait.OOMKilled {
		err = &OomError{
			Name: step.Metadata.Name,
			Code: wait.ExitCode,
		}
	} else if wait.ExitCode != 0 {
		err = &ExitError{
			Name: step.Metadata.Name,
			Code: wait.ExitCode,
		}
	}

	if r.hook.AfterEach != nil {
		state := snapshot(r, step, wait)
		if err := r.hook.AfterEach(state); err != nil {
			return err
		}
	}

	if step.IgnoreErr {
		return nil
	}
	return err
}

// helper function exports a single file or folder.
func stream(state *State, rc io.ReadCloser) error {
	defer rc.Close()

	w := newWriter(state)
	io.Copy(w, rc)

	if state.hook.GotLogs != nil {
		return state.hook.GotLogs(state, w.lines)
	}
	return nil
}
