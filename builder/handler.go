package builder

import (
	"io"

	"github.com/drone/drone-cli/cluster"
)

type Handler interface {
	Run(r *Request, rw ResultWriter) error
	Cancel() error
}

type BatchHandler struct {
	cluster   cluster.Cluster
	container *cluster.Container
}

func (b *BatchHandler) Run(r *Request, rw ResultWriter) error {
	b.container = &cluster.Container{
		Image:       r.Config.Clone.Image,
		Pull:        r.Config.Clone.Pull,
		Privileged:  r.Config.Clone.Privileged,
		Env:         r.Config.Clone.Environment,
		Volumes:     r.Config.Clone.Volumes,
		NetworkMode: r.Config.Clone.Net,
	}
	err := b.cluster.Create(b.container)
	if err != nil {
		return err
	}
	err = b.cluster.Start(b.container)
	if err != nil {
		return err
	}
	rc, err := b.cluster.Logs(b.container)
	if err != nil {
		return err
	}
	io.Copy(rw, rc)
	rc.Close()
	state, err := b.cluster.State(b.container)
	if err != nil {
		return err
	}
	rw.WriteExitCode(state.ExitCode)
	if state.ExitCode != 0 {
		break
	}
	return nil
}

func (b *BatchHandler) Cancel() error {
	b.cluster.Stop(b.container)
	return b.cluster.Remove(b.container)
}

type ServiceHandler struct {
	cluster   cluster.Cluster
	container *cluster.Container
}

func (s *ServiceHandler) Cancel() error {
	s.cluster.Stop(s.container)
	return s.cluster.Remove(s.container)
}

func (s *ServiceHandler) Run(r *Request, rw ResultWriter) error {
	s.container = &cluster.Container{
		Image:       r.Config.Clone.Image,
		Pull:        r.Config.Clone.Pull,
		Privileged:  r.Config.Clone.Privileged,
		Env:         r.Config.Clone.Environment,
		Volumes:     r.Config.Clone.Volumes,
		NetworkMode: r.Config.Clone.Net,
	}
	err := s.cluster.Create(s.container)
	if err != nil {
		return err
	}
	return s.cluster.Start(s.container)
}

//
//
//

// type BatchHandler interface{}
// type ServiceHandler interface{}
// type BuildHandler interface{}

// build.Handle(builder.BatchHandler(env, step))
// build.Handle(builder.ServiceHandler(env, step))
// build.Handle(builder.ServiceHandler(env, step))
// build.Handle(builder.ServiceHandler(env, step))
// build.Handle(builder.BuildHandler(env, step))
// build.Run()

// deploy.Handle(builder.BatchHandler(step))
// deploy.Handle(builder.BatchHandler(step))
// deploy.Handle(builder.BatchHandler(step))
// deploy.Handle(builder.BatchHandler(step))
// deploy.Run()

// report.Handle(builder.BatchHandler(step))
// report.Handle(builder.BatchHandler(step))
// report.Handle(builder.BatchHandler(step))
// report.Handle(builder.BatchHandler(step))
// report.Run()
