package builder

import (
	"io"

	"github.com/drone/drone-cli/cluster"
)

// type Builder interface {
// 	Register(*cluster.Container) error
// 	Run(ResultWriter) error
// 	Cancel() error
// }

type Builder struct {
	cluster    cluster.Cluster
	containers []*cluster.Container
}

func (b *Builder) Register(c *cluster.Container) error {
	err := b.cluster.Create(c)
	if err != nil {
		return err
	}
	b.containers = append(b.containers, c)
	return nil
}

func (b *Builder) Build(rw ResultWriter) error {
	for _, c := range b.containers {
		b.cluster.Start(c)
		if c.Detach {
			continue
		}

		r, err := b.cluster.Logs(c)
		if err != nil {
			return err
		}
		io.Copy(rw, r)
		r.Close()
		state, err := b.cluster.State(c)
		if err != nil {
			return err
		}

		rw.WriteExitCode(state.ExitCode)
		if state.ExitCode != 0 {
			break
		}
	}
	return nil
}

func (b *Builder) Cancel() {
	for _, c := range b.containers {
		b.cluster.Stop(c)
		b.cluster.Remove(c)
	}
}
