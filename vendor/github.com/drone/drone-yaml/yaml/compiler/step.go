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

package compiler

import (
	"strings"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml/yaml"
	"github.com/drone/drone-yaml/yaml/compiler/image"
	"github.com/drone/drone-yaml/yaml/compiler/internal/rand"
)

func createStep(spec *engine.Spec, src *yaml.Container) *engine.Step {
	dst := &engine.Step{
		Metadata: engine.Metadata{
			UID:       rand.String(),
			Name:      src.Name,
			Namespace: spec.Metadata.Namespace,
			Labels: map[string]string{
				"io.drone.step.name": src.Name,
			},
		},
		Detach:    src.Detach,
		DependsOn: src.DependsOn,
		Devices:   nil,
		Docker: &engine.DockerStep{
			Args:       src.Command,
			Command:    src.Entrypoint,
			DNS:        src.DNS,
			DNSSearch:  src.DNSSearch,
			ExtraHosts: src.ExtraHosts,
			Image:      image.Expand(src.Image),
			Networks:   nil, // set in compiler.go
			Ports:      toPorts(src),
			Privileged: src.Privileged,
			PullPolicy: toPullPolicy(src),
		},
		Envs:         map[string]string{},
		Files:        nil, // set below
		IgnoreErr:    toIgnoreErr(src),
		IgnoreStderr: false,
		IgnoreStdout: false,
		Resources:    toResources(src),
		RunPolicy:    toRunPolicy(src),
		Secrets:      nil, // set below
		Volumes:      nil, // set below
		WorkingDir:   "",  // set in compiler.go
	}

	// if the user is running a service container with
	// no custom commands, we should revert back to the
	// user-defined working directory, which may be empty.
	if dst.Detach && len(src.Commands) == 0 {
		dst.WorkingDir = src.WorkingDir
	}

	// appends the volumes to the container def.
	for _, vol := range src.Volumes {
		// the user should never be able to directly
		// mount the docker socket. This should be
		// restricted by the linter, but we place this
		// check here just to be safe.
		if vol.Name == "_docker_socket" {
			continue
		}
		mount := &engine.VolumeMount{
			Name: vol.Name,
			Path: vol.MountPath,
		}
		dst.Volumes = append(dst.Volumes, mount)
	}

	// appends the environment variables to the
	// container definition.
	for key, value := range src.Environment {
		// fix https://github.com/drone/drone-yaml/issues/13
		if value == nil {
			continue
		}
		if value.Secret != "" {
			sec := &engine.SecretVar{
				Name: value.Secret,
				Env:  key,
			}
			dst.Secrets = append(dst.Secrets, sec)
		} else {
			dst.Envs[key] = value.Value
		}
	}

	// appends the settings variables to the
	// container definition.
	for key, value := range src.Settings {
		// fix https://github.com/drone/drone-yaml/issues/13
		if value == nil {
			continue
		}
		// all settings are passed to the plugin env
		// variables, prefixed with PLUGIN_
		key = "PLUGIN_" + strings.ToUpper(key)

		// if the setting parameter is sources from the
		// secret we create a secret enviornment variable.
		if value.Secret != "" {
			sec := &engine.SecretVar{
				Name: value.Secret,
				Env:  key,
			}
			dst.Secrets = append(dst.Secrets, sec)
		} else {
			// else if the setting parameter is opaque
			// we inject as a string-encoded environment
			// variable.
			dst.Envs[key] = encode(value.Value)
		}
	}

	// if the step specifies shell commands we generate a
	// script. The script is copied to the container at
	// runtime (or mounted as a config map) and then executed
	// as the entrypoint.
	if len(src.Commands) > 0 {
		switch spec.Platform.OS {
		case "windows":
			setupScriptWin(spec, dst, src)
		default:
			setupScript(spec, dst, src)
		}
	}

	return dst
}

func createBuildStep(spec *engine.Spec, src *yaml.Container) *engine.Step {
	dst := &engine.Step{
		Metadata: engine.Metadata{
			UID:       rand.String(),
			Name:      src.Name,
			Namespace: spec.Metadata.Namespace,
			Labels: map[string]string{
				"io.drone.step.name": src.Name,
			},
		},
		Detach:    src.Detach,
		DependsOn: src.DependsOn,
		Devices:   nil,
		Docker: &engine.DockerStep{
			Args:       []string{"--build"},
			DNS:        src.DNS,
			Image:      "drone/docker",
			PullPolicy: engine.PullIfNotExists,
		},
		Envs:         map[string]string{},
		IgnoreErr:    toIgnoreErr(src),
		IgnoreStderr: false,
		IgnoreStdout: false,
		Resources:    toResources(src),
		RunPolicy:    toRunPolicy(src),
	}

	// if v := src.Build.Args; len(v) > 0 {
	// 	dst.Envs["DOCKER_BUILD_ARGS"] = strings.Join(v, ",")
	// }
	if v := src.Build.CacheFrom; len(v) > 0 {
		dst.Envs["DOCKER_BUILD_CACHE_FROM"] = strings.Join(v, ",")
	}
	// if len(src.Build.Labels) > 0 {
	// 	dst.Envs["DOCKER_BUILD_LABELS"] = strings.Join(v, ",")
	// }
	if v := src.Build.Dockerfile; v != "" {
		dst.Envs["DOCKER_BUILD_DOCKERFILE"] = v

	}
	if v := src.Build.Context; v != "" {
		dst.Envs["DOCKER_BUILD_CONTEXT"] = v
	}
	if v := src.Build.Image; v != "" {
		alias := image.Trim(v) + ":" + dst.Metadata.UID
		dst.Envs["DOCKER_BUILD_IMAGE"] = image.Expand(v)
		dst.Envs["DOCKER_BUILD_IMAGE_ALIAS"] = image.Expand(alias)
	}

	dst.Volumes = append(dst.Volumes, &engine.VolumeMount{
		Name: "_docker_socket",
		Path: "/var/run/docker.sock",
	})

	return dst
}
