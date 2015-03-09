package builder

import (
	"io"

	"github.com/drone/drone-cli/common"
	"github.com/drone/drone-cli/engine"
)

type Handler interface {
	Run(ResponseWriter) error
	Cancel()
}

type handler struct {
	engine    engine.Engine
	container *engine.Container
	detached  bool
}

func (h *handler) Run(rw ResponseWriter) error {
	err := h.engine.Create(h.container)
	if err != nil {
		return err
	}
	err = h.engine.Start(h.container)
	if err != nil {
		return err
	}
	if h.detached {
		return nil
	}
	rc, err := h.engine.Logs(h.container)
	if err != nil {
		return err
	}
	io.Copy(rw, rc)
	rc.Close()
	state, err := h.engine.State(h.container)
	if err != nil {
		return err
	}
	rw.WriteExitCode(state.ExitCode)
	return nil
}

func (h *handler) Cancel() {
	h.engine.Stop(h.container)
	h.engine.Remove(h.container)
}

func BatchHandler(step *common.Step) Handler {
	container := &engine.Container{
		Image:       step.Image,
		Pull:        step.Pull,
		Privileged:  step.Privileged,
		Env:         step.Environment,
		Volumes:     step.Volumes,
		NetworkMode: step.Net,
		Cmd:         []string{}, // TODO need to encode data
	}
	return &handler{container: container}
}

func ServiceHandler(step *common.Step) Handler {
	return &handler{
		detached: true,
		container: &engine.Container{
			Image:       step.Image,
			Pull:        step.Pull,
			Privileged:  step.Privileged,
			Env:         step.Environment,
			Volumes:     step.Volumes,
			NetworkMode: step.Net,
		}}
}

func BuildHandler(step *common.Step) Handler {
	container := &engine.Container{
		Image:       step.Image,
		Pull:        step.Pull,
		Privileged:  step.Privileged,
		Env:         step.Environment,
		Volumes:     step.Volumes,
		NetworkMode: step.Net,
		Entrypoint:  []string{"/bin/bash"},
		Cmd:         []string{"/drone/bin/build.sh"},
	}
	return &handler{container: container}
}

func SetupHandler(step *common.Step) Handler {
	container := &engine.Container{
		Image: "plugin/drone-build",
		Cmd:   []string{"/drone/bin/build.sh"}, // TODO need to encode data. Encode BUILD Args
	}
	return &handler{container: container}
}

type service struct {
	engine    engine.Engine
	container *engine.Container
}

func (s *service) Run(rw ResponseWriter, r *Request) error {
	err := s.engine.Create(s.container)
	if err != nil {
		return err
	}
	return s.engine.Start(s.container)
}

type batch struct {
	engine    engine.Engine
	container *engine.Container
}

func (b *batch) Run(rw ResponseWriter, r *Request) error {
	err := b.engine.Create(b.container)
	if err != nil {
		return err
	}
	err = b.engine.Start(b.container)
	if err != nil {
		return err
	}
	rc, err := b.engine.Logs(b.container)
	if err != nil {
		return err
	}
	io.Copy(rw, rc)
	rc.Close()
	state, err := b.engine.State(b.container)
	if err != nil {
		return err
	}
	rw.WriteExitCode(state.ExitCode)
	return nil
}
