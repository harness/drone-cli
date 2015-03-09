package docker

import (
	"io"
	"io/ioutil"

	"github.com/drone/drone-cli/engine"
	"github.com/samalba/dockerclient"
)

func create(client dockerclient.Client, container *engine.Container) error {
	conf := dockerclient.ContainerConfig{
		Image:      container.Image,
		Cmd:        container.Cmd,
		Env:        container.Env,
		Entrypoint: container.Entrypoint,
		WorkingDir: container.WorkingDir,
	}
	conf.Volumes = map[string]struct{}{}
	for _, path := range container.Volumes {
		conf.Volumes[path] = struct{}{}
	}

	id, err := client.CreateContainer(&conf, "")
	if err != nil {
		return err
	}
	container.ID = id
	return nil
}

func start(client dockerclient.Client, ambassador, container *engine.Container) error {
	config := dockerclient.HostConfig{
		Privileged:  container.Privileged,
		NetworkMode: container.NetworkMode,
	}
	if len(container.NetworkMode) == 0 && ambassador != nil {
		config.NetworkMode = "container:" + ambassador.ID
	}
	if ambassador != nil {
		config.VolumesFrom = []string{"container:" + ambassador.ID}
	}
	return client.StartContainer(container.ID, &config)
}

func stop(client dockerclient.Client, container *engine.Container) error {
	client.StopContainer(container.ID, 5)
	client.KillContainer(container.ID, "SIGKILL")
	return nil
}

func remove(client dockerclient.Client, container *engine.Container) error {
	return client.RemoveContainer(container.ID, true, false)
}

func state(client dockerclient.Client, container *engine.Container) (*engine.State, error) {
	info, err := client.InspectContainer(container.ID)
	if err != nil {
		return nil, err
	}
	return &engine.State{
		Running:  info.State.Running,
		Pid:      info.State.Pid,
		ExitCode: info.State.ExitCode,
		Started:  info.State.StartedAt.UTC().Unix(),
		Finished: info.State.FinishedAt.UTC().Unix(),
	}, nil
}

func logs(client dockerclient.Client, container *engine.Container) (io.ReadCloser, error) {
	opts := dockerclient.LogOptions{
		Follow:     true,
		Stderr:     true,
		Stdout:     true,
		Timestamps: true,
	}
	return client.ContainerLogs(container.ID, &opts)
}

func wait(client dockerclient.Client, container *engine.Container) error {
	src, err := logs(client, container)
	if err != nil {
		return err
	}
	defer src.Close()
	_, err = io.Copy(ioutil.Discard, src)
	return nil
}
