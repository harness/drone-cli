package runner

import (
	"io"
	"sync"

	"github.com/drone/drone-cli/common"
	"github.com/drone/drone-cli/common/uuid"
	"github.com/samalba/dockerclient"
)

type Builder struct {
	client dockerclient.Client
	config *common.Config

	ambassador *Container
	containers []*Container
	notifiers  []*Container

	sync.RWMutex
}

func (b *Builder) Run(req *Request, resp ResponseWriter) error {
	err := b.ambassador.Create()
	if err != nil {
		return err
	}
	err = b.ambassador.Start()
	if err != nil {
		return err
	}

	for _, container := range b.containers {
		err = container.Create()
		if err != nil {
			return err
		}
		err = container.Start()
		if err != nil {
			return err
		}

		if container.Detached {
			continue
		}

		r, err := container.Logs()
		if err != nil {
			return err
		}
		io.Copy(resp, r)
		r.Close()
		info, err := container.Inspect()
		if err != nil {
			return err
		}

		if info.State.Running != false {
			println("ERROR: container still running")
		}

		resp.WriteExitCode(info.State.ExitCode)
		if info.State.ExitCode != 0 {
			break
		}
	}

	for _, container := range b.notifiers {
		container.Create()
		container.Start()
		container.Logs()
	}

	return nil
}

func (b *Builder) Setup(req *Request) error {
	uid := uuid.CreateUUID()
	net := "container:" + uid

	// creates an ambassador container that can be
	// used for the shared network and mounted volume.
	b.ambassador = &Container{
		Name:       uid,
		Image:      "busybox",
		Volumes:    []string{"/drone"},
		Entrypoint: []string{"/bin/sleep", "1d"},
	}

	// initializes our build environment
	b.containers = append(b.containers, &Container{
		Image:       b.config.Init.Image,
		Volumes:     []string{"/drone"},
		Env:         b.config.Init.Environment,
		Privileged:  b.config.Init.Privileged,
		NetworkMode: net,
		Pull:        b.config.Init.Pull,
		Cmd:         EncodeParams(req, req.Config.Build.Config),
	})

	// clones the repository
	b.containers = append(b.containers, &Container{
		Image:       b.config.Clone.Image,
		VolumesFrom: []string{b.ambassador.Name},
		Env:         b.config.Clone.Environment,
		Privileged:  b.config.Clone.Privileged,
		NetworkMode: net,
		Pull:        b.config.Clone.Pull,
		Cmd:         EncodeParams(req, req.Config.Clone.Config),
	})

	// service containers
	for _, service := range req.Config.Compose {
		b.containers = append(b.containers, &Container{
			Image:       service.Image,
			Pull:        service.Pull,
			Env:         service.Environment,
			Privileged:  service.Privileged,
			NetworkMode: net,
			Detached:    true,
		})
	}

	// runs our build and tests
	b.containers = append(b.containers, &Container{
		Image:       b.config.Build.Image,
		VolumesFrom: []string{b.ambassador.Name},
		Env:         b.config.Build.Environment,
		Privileged:  b.config.Build.Privileged,
		NetworkMode: net,
		Pull:        b.config.Build.Pull,
		Cmd:         EncodeParams(req, req.Config.Clone.Config),
	})

	// deploy containers
	for _, deploy := range req.Config.Deploy {
		b.containers = append(b.containers, &Container{
			Image:       deploy.Image,
			Pull:        deploy.Pull,
			Env:         deploy.Environment,
			Privileged:  deploy.Privileged,
			NetworkMode: net,
			VolumesFrom: []string{b.ambassador.Name},
			Cmd:         EncodeParams(req, deploy.Config),
		})
	}

	// notify containers
	for _, notify := range req.Config.Deploy {
		b.notifiers = append(b.notifiers, &Container{
			Image:      notify.Image,
			Pull:       notify.Pull,
			Env:        notify.Environment,
			Privileged: notify.Privileged,
			Cmd:        EncodeParams(req, notify.Config),
		})
	}

	return nil
}

func (b *Builder) Teardown() {
	for _, container := range b.notifiers {
		container.Stop()
		container.Kill()
		container.Remove()
	}
	for _, container := range b.containers {
		container.Stop()
		container.Kill()
		container.Remove()
	}
	b.ambassador.Stop()
	b.ambassador.Kill()
	b.ambassador.RemoveVolume()
}

func (b *Builder) Logs(w io.Writer) error {
	for _, container := range b.containers {
		if container.Detached {
			continue
		}
		rc, err := container.Logs()
		if err != nil {
			return err
		}
		io.Copy(w, rc)
	}
	return nil
}
