package builder

import "github.com/drone/drone-cli/runner/builder/runner"

type Builder struct {
	environ  runner.Environment
	commands []*runner.Command
}

func (b *Builder) Push(cmd *runner.Command) error {
	err := b.environ.Create(cmd)
	if err != nil {
		return err
	}
	b.commands = append(b.commands, cmd)
	return nil
}

func (b *Builder) Run() error {
	for _, cmd := range b.commands {
		b.environ.Start(cmd)
		if cmd.Detach {
			continue
		}
	}
	return nil
}

func (b *Builder) Cancel() {
	for _, cmd := range b.commands {
		b.environ.Stop(cmd)
		b.environ.Remove(cmd)
	}
}
