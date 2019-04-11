// Copyright 2019 Drone IO, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package docker

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"time"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-runtime/engine/docker/auth"
	"github.com/drone/drone-runtime/engine/docker/stdcopy"

	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/network"
	"docker.io/go-docker/api/types/volume"
)

type dockerEngine struct {
	client docker.APIClient
}

// NewEnv returns a new Engine from the environment.
func NewEnv() (engine.Engine, error) {
	cli, err := docker.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return New(cli), nil
}

// NewEnvBackoff returns a new Engine from the environment. If
// the connection fails or cannot be reached, the system attempts
// to reconnect N times before returning an error.
func NewEnvBackoff(retries int, wait time.Duration) (engine.Engine, error) {
	var cli *docker.Client
	var err error
	for i := 0; i < retries; i++ {
		if i > 0 {
			time.Sleep(wait)
		}
		cli, err = docker.NewEnvClient()
		if err == nil {
			continue
		}
		_, err = cli.Ping(context.Background())
		if err != nil {
			cli.Close()
			continue
		}
		break
	}
	if err != nil {
		return nil, err
	}
	return New(cli), nil
}

// New returns a new Engine using the Docker API Client.
func New(client docker.APIClient) engine.Engine {
	return &dockerEngine{
		client: client,
	}
}

func (e *dockerEngine) Setup(ctx context.Context, spec *engine.Spec) error {
	if spec.Docker != nil {
		// creates the default temporary (local) volumes
		// that are mounted into each container step.
		for _, vol := range spec.Docker.Volumes {
			if vol.EmptyDir == nil {
				continue
			}

			_, err := e.client.VolumeCreate(ctx, volume.VolumesCreateBody{
				Name:   vol.Metadata.UID,
				Driver: "local",
				Labels: spec.Metadata.Labels,
			})
			if err != nil {
				return err
			}
		}
	}

	// creates the default pod network. All containers
	// defined in the pipeline are attached to this network.
	driver := "bridge"
	if spec.Platform.OS == "windows" {
		driver = "nat"
	}
	_, err := e.client.NetworkCreate(ctx, spec.Metadata.UID, types.NetworkCreate{
		Driver: driver,
		Labels: spec.Metadata.Labels,
	})

	return err
}

func (e *dockerEngine) Create(ctx context.Context, spec *engine.Spec, step *engine.Step) error {
	if step.Docker == nil {
		return errors.New("engine: missing docker configuration")
	}

	// parse the docker image name. We need to extract the
	// image domain name and match to registry credentials
	// stored in the .docker/config.json object.
	_, domain, latest, err := parseImage(step.Docker.Image)
	if err != nil {
		return err
	}

	// create pull options with encoded authorization credentials.
	pullopts := types.ImagePullOptions{}
	auths, ok := engine.LookupAuth(spec, domain)
	if ok {
		pullopts.RegistryAuth = auth.Encode(auths.Username, auths.Password)
	}

	// automatically pull the latest version of the image if requested
	// by the process configuration.
	if step.Docker.PullPolicy == engine.PullAlways ||
		(step.Docker.PullPolicy == engine.PullDefault && latest) {
		// TODO(bradrydzewski) implement the PullDefault strategy to pull
		// the image if the tag is :latest
		rc, perr := e.client.ImagePull(ctx, step.Docker.Image, pullopts)
		if perr == nil {
			io.Copy(ioutil.Discard, rc)
			rc.Close()
		}
		if perr != nil {
			return perr
		}
	}

	_, err = e.client.ContainerCreate(ctx,
		toConfig(spec, step),
		toHostConfig(spec, step),
		toNetConfig(spec, step),
		step.Metadata.UID,
	)

	// automatically pull and try to re-create the image if the
	// failure is caused because the image does not exist.
	if docker.IsErrImageNotFound(err) && step.Docker.PullPolicy != engine.PullNever {
		rc, perr := e.client.ImagePull(ctx, step.Docker.Image, pullopts)
		if perr != nil {
			return perr
		}
		io.Copy(ioutil.Discard, rc)
		rc.Close()

		// once the image is successfully pulled we attempt to
		// re-create the container.
		_, err = e.client.ContainerCreate(ctx,
			toConfig(spec, step),
			toHostConfig(spec, step),
			toNetConfig(spec, step),
			step.Metadata.UID,
		)
	}
	if err != nil {
		return err
	}

	copyOpts := types.CopyToContainerOptions{}
	copyOpts.AllowOverwriteDirWithFile = false
	for _, mount := range step.Files {
		file, ok := engine.LookupFile(spec, mount.Name)
		if !ok {
			continue
		}
		tar := createTarfile(file, mount)

		// TODO(bradrydzewski) this path is probably different on windows.
		err := e.client.CopyToContainer(ctx, step.Metadata.UID, "/", bytes.NewReader(tar), copyOpts)
		if err != nil {
			return err
		}
	}

	for _, net := range step.Docker.Networks {
		err = e.client.NetworkConnect(ctx, net, step.Metadata.UID, &network.EndpointSettings{
			Aliases: []string{net},
		})
		if err != nil {
			return nil
		}
	}

	return nil
}

func (e *dockerEngine) Start(ctx context.Context, spec *engine.Spec, step *engine.Step) error {
	return e.client.ContainerStart(ctx, step.Metadata.UID, types.ContainerStartOptions{})
}

func (e *dockerEngine) Wait(ctx context.Context, spec *engine.Spec, step *engine.Step) (*engine.State, error) {
	wait, errc := e.client.ContainerWait(ctx, step.Metadata.UID, "")
	select {
	case <-wait:
	case <-errc:
	}

	info, err := e.client.ContainerInspect(ctx, step.Metadata.UID)
	if err != nil {
		return nil, err
	}
	if info.State.Running {
		// TODO(bradrydewski) if the state is still running
		// we should call wait again.
	}

	return &engine.State{
		Exited:    true,
		ExitCode:  info.State.ExitCode,
		OOMKilled: info.State.OOMKilled,
	}, nil
}

func (e *dockerEngine) Tail(ctx context.Context, spec *engine.Spec, step *engine.Step) (io.ReadCloser, error) {
	opts := types.ContainerLogsOptions{
		Follow:     true,
		ShowStdout: true,
		ShowStderr: true,
		Details:    false,
		Timestamps: false,
	}

	logs, err := e.client.ContainerLogs(ctx, step.Metadata.UID, opts)
	if err != nil {
		return nil, err
	}
	rc, wc := io.Pipe()

	go func() {
		stdcopy.StdCopy(wc, wc, logs)
		logs.Close()
		wc.Close()
		rc.Close()
	}()
	return rc, nil
}

func (e *dockerEngine) Destroy(ctx context.Context, spec *engine.Spec) error {
	removeOpts := types.ContainerRemoveOptions{
		Force:         true,
		RemoveLinks:   false,
		RemoveVolumes: true,
	}

	// stop all containers
	for _, step := range spec.Steps {
		e.client.ContainerKill(ctx, step.Metadata.UID, "9")
	}

	// cleanup all containers
	for _, step := range spec.Steps {
		e.client.ContainerRemove(ctx, step.Metadata.UID, removeOpts)
	}

	// cleanup all volumes
	if spec.Docker != nil {
		for _, vol := range spec.Docker.Volumes {
			if vol.EmptyDir == nil {
				continue
			}
			// tempfs volumes do not have a volume entry,
			// and therefore do not require removal.
			if vol.EmptyDir.Medium == "memory" {
				continue
			}
			e.client.VolumeRemove(ctx, vol.Metadata.UID, true)
		}
	}

	// cleanup the network
	e.client.NetworkRemove(ctx, spec.Metadata.UID)

	// notice that we never collect or return any errors.
	// this is because we silently ignore cleanup failures
	// and instead ask the system admin to periodically run
	// `docker prune` commands.
	return nil
}
