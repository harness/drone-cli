package builder

import (
	"io"

	"github.com/drone/drone-cli/runner/builder/runner"
)

type Notifier struct {
	cluster  runner.Cluster
	commands []*runner.Container
}

func (n *Notifier) Register(cmd *runner.Container) error {
	err := n.cluster.Create(cmd)
	if err != nil {
		return err
	}
	n.commands = append(n.commands, cmd)
	return nil
}

func (n *Notifier) Run(rw *runner.Result) error {
	for _, cmd := range n.commands {
		n.cluster.Start(cmd)
		if cmd.Detach {
			continue
		}

		r, err := n.cluster.Logs(cmd)
		if err != nil {
			return err
		}
		io.Copy(rw, r)
		r.Close()
		state, err := n.cluster.State(cmd)
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

func (n *Notifier) Cancel() {
	for _, cmd := range n.commands {
		n.cluster.Stop(cmd)
		n.cluster.Remove(cmd)
	}
}
