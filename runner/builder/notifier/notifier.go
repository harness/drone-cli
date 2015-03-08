package builder

import "github.com/samalba/dockerclient"

type Notifier struct {
	environs *Environment
	commands []*Command
}

func (n *Notifier) Register(c *Command) error {
	err := create(c, n.environs.client)
	if err != nil {
		return err
	}
	n.commands = append(n.commands, c)
	return nil
}

func (n *Notifier) Run() error {
	for _, c := range n.commands {
		n.environs.client.StartContainer(c.ID, &dockerclient.HostConfig{})
		wait(c, n.environs.client)
	}
	return nil
}

func (n *Notifier) Destroy() {
	for _, c := range n.commands {
		n.environs.client.StopContainer(c.ID, 5)
		n.environs.client.KillContainer(c.ID, "SIGKILL")
		n.environs.client.RemoveContainer(c.ID, true, false)
	}
}

type Runner interface {
	Add(*Command) error
	Run(*Command) error
	Setup() error
	Teardown() error
}
