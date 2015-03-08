package builder

import (
	"io"
	//"io/ioutil"

	"github.com/samalba/dockerclient"
)

// An Runner represents a build exeuction environment used
// to execute a single build.
type Runner struct {
	client      dockerclient.Client
	environment *Environment
	containers  []*Command
}

func (r *Runner) Push(c *Command) error {
	err := create(c, r.client)
	if err != nil {
		return err
	}
	r.containers = append(r.containers, c)
	return nil
}

func (r *Runner) Run(w ResponseWriter) error {
	var err error
	for _, c := range r.containers {

		if c.Defer {
			continue
		}

		config := dockerclient.HostConfig{
			Privileged:  c.Privileged,
			NetworkMode: r.environment.String(),
			VolumesFrom: []string{r.environment.String()},
		}
		err = r.client.StartContainer(c.ID, &config)
		if err != nil {
			break
		}

		if c.Detach {
			continue
		}

		rc, err := logs(c, r.client)
		if err != nil {
			break
		}
		io.Copy(w, rc)
		rc.Close()
		info, err := r.client.InspectContainer(c.ID)
		if err != nil {
			break
		}

		w.WriteExitCode(info.State.ExitCode)
		if info.State.ExitCode != 0 {
			break
		}
	}

	for _, c := range r.containers {
		if !c.Defer {
			continue
		}
		r.client.StartContainer(c.ID, &dockerclient.HostConfig{})
		wait(c, r.client)
	}

	return err
}

func Build() {

}

func Deploy() {

}

func (b *Runner) Notify() error {
	for _, c := range b.containers {
		conf := dockerclient.HostConfig{}
		err := b.client.StartContainer(c.ID, &conf)
		if err != nil {
			return err
		}
		wait(c, b.client)
	}
	return nil
}

func Setup() {

}

func Teardown() {

}
