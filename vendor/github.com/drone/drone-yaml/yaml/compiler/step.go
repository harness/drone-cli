package compiler

import (
	"strings"

	"github.com/drone/drone-runtime/engine"
	"github.com/drone/drone-yaml/yaml"
	"github.com/drone/drone-yaml/yaml/compiler/image"
)

func createStep(spec *engine.Spec, src *yaml.Container) *engine.Step {
	dst := &engine.Step{
		Metadata: engine.Metadata{
			UID:       random(),
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
		mount := &engine.VolumeMount{
			Name: vol.Name,
			Path: vol.MountPath,
		}
		dst.Volumes = append(dst.Volumes, mount)
	}

	// appends the environment variables to the
	// container definition.
	for key, value := range src.Environment {
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
