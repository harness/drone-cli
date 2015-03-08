package preparer

import (
	"github.com/drone/drone-cli/builder"
	"github.com/drone/drone-cli/cluster"
)

func PrepClone(b builder.Builder, r builder.Request) error {
	if len(r.Config.Clone.Image) == 0 {
		return nil
	}
	return b.Register(&cluster.Container{
		Image:       r.Config.Clone.Image,
		Pull:        r.Config.Clone.Pull,
		Privileged:  r.Config.Clone.Privileged,
		Env:         r.Config.Clone.Environment,
		Volumes:     r.Config.Clone.Volumes,
		NetworkMode: r.Config.Clone.Net,
	})
}

func PrepBuild(b builder.Builder, r builder.Request) error {
	if len(r.Config.Build.Image) == 0 {
		return nil
	}
	return b.Register(&cluster.Container{
		Image:       r.Config.Build.Image,
		Pull:        r.Config.Build.Pull,
		Privileged:  r.Config.Build.Privileged,
		Env:         r.Config.Build.Environment,
		Volumes:     r.Config.Build.Volumes,
		NetworkMode: r.Config.Build.Net,
		WorkingDir:  r.Clone.Dir,
		Entrypoint:  []string{"/bin/bash"},
		Cmd:         []string{"/drone/bin/build.sh"},
	})
}

func PrepCompose(b builder.Builder, r builder.Request) error {
	if len(r.Config.Compose) == 0 {
		return nil
	}
	var err error
	for _, config := range r.Config.Compose {
		err = b.Register(&cluster.Container{
			Image:       config.Image,
			Pull:        config.Pull,
			Privileged:  config.Privileged,
			Env:         config.Environment,
			Volumes:     config.Volumes,
			NetworkMode: config.Net,
			Detach:      true,
		})
		if err != nil {
			break
		}
	}
	return err
}

func PrepPublish(b builder.Builder, r builder.Request) error {
	if len(r.Config.Publish) == 0 {
		return nil
	}
	var err error
	for _, config := range r.Config.Publish {
		err = b.Register(&cluster.Container{
			Image:       config.Image,
			Pull:        config.Pull,
			Privileged:  config.Privileged,
			Env:         config.Environment,
			Volumes:     config.Volumes,
			NetworkMode: config.Net,
		})
		if err != nil {
			break
		}
	}
	return err
}

func PrepDeploy(b builder.Builder, r builder.Request) error {
	if len(r.Config.Deploy) == 0 {
		return nil
	}
	var err error
	for _, config := range r.Config.Deploy {
		err = b.Register(&cluster.Container{
			Image:       config.Image,
			Pull:        config.Pull,
			Privileged:  config.Privileged,
			Env:         config.Environment,
			Volumes:     config.Volumes,
			NetworkMode: config.Net,
		})
		if err != nil {
			break
		}
	}
	return err
}

func PrepNotify(b builder.Builder, r builder.Request) error {
	if len(r.Config.Notify) == 0 {
		return nil
	}
	var err error
	for _, config := range r.Config.Notify {
		err = b.Register(&cluster.Container{
			Image:       config.Image,
			Pull:        config.Pull,
			Privileged:  config.Privileged,
			Env:         config.Environment,
			Volumes:     config.Volumes,
			NetworkMode: config.Net,
		})
		if err != nil {
			break
		}
	}
	return err
}
