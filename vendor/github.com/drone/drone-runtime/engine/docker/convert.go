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
	"strings"

	"github.com/drone/drone-runtime/engine"

	"docker.io/go-docker/api/types/container"
	"docker.io/go-docker/api/types/mount"
	"docker.io/go-docker/api/types/network"
)

// returns a container configuration.
func toConfig(spec *engine.Spec, step *engine.Step) *container.Config {
	config := &container.Config{
		Image:        step.Docker.Image,
		Labels:       step.Metadata.Labels,
		WorkingDir:   step.WorkingDir,
		User:         step.Docker.User,
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		OpenStdin:    false,
		StdinOnce:    false,
		ArgsEscaped:  false,
	}

	if len(step.Envs) != 0 {
		config.Env = toEnv(step.Envs)
	}
	for _, sec := range step.Secrets {
		secret, ok := engine.LookupSecret(spec, sec)
		if ok {
			config.Env = append(config.Env, sec.Env+"="+secret.Data)
		}
	}
	if len(step.Docker.Args) != 0 {
		config.Cmd = step.Docker.Args
	}
	if len(step.Docker.Command) != 0 {
		config.Entrypoint = step.Docker.Command
	}

	// NOTE it appears this is no longer required,
	// however this could cause incompatibility with
	// certain docker versions.
	//
	//   if len(step.Volumes) != 0 {
	// 	    config.Volumes = toVolumeSet(spec, step)
	//   }
	return config
}

// returns a container host configuration.
func toHostConfig(spec *engine.Spec, step *engine.Step) *container.HostConfig {
	config := &container.HostConfig{
		LogConfig: container.LogConfig{
			Type: "json-file",
		},
		Privileged: step.Docker.Privileged,
		// TODO(bradrydzewski) set ShmSize
	}
	// windows does not support privileged so we hard-code
	// this value to false.
	if spec.Platform.OS == "windows" {
		config.Privileged = false
	}
	if len(step.Docker.DNS) > 0 {
		config.DNS = step.Docker.DNS
	}
	if len(step.Docker.DNSSearch) > 0 {
		config.DNSSearch = step.Docker.DNSSearch
	}
	if len(step.Docker.ExtraHosts) > 0 {
		config.ExtraHosts = step.Docker.ExtraHosts
	}
	if step.Resources != nil {
		config.Resources = container.Resources{}
		if limits := step.Resources.Limits; limits != nil {
			config.Resources.Memory = limits.Memory
			// TODO(bradrydewski) set config.Resources.CPUPercent

			// IMPORTANT docker and kubernetes use
			// different units of measure for cpu limits.
			// we need to figure out how to convert from
			// the kubernetes unit of measure to the docker
			// unit of measure.
		}
	}

	if len(step.Volumes) != 0 {
		config.Devices = toDeviceSlice(spec, step)
		config.Binds = toVolumeSlice(spec, step)
		config.Mounts = toVolumeMounts(spec, step)
	}
	return config
}

// helper function returns the container network configuration.
func toNetConfig(spec *engine.Spec, proc *engine.Step) *network.NetworkingConfig {
	endpoints := map[string]*network.EndpointSettings{}
	endpoints[spec.Metadata.UID] = &network.EndpointSettings{
		NetworkID: spec.Metadata.UID,
		Aliases:   []string{proc.Metadata.Name},
	}
	return &network.NetworkingConfig{
		EndpointsConfig: endpoints,
	}
}

// helper function that converts a slice of device paths to a slice of
// container.DeviceMapping.
func toDeviceSlice(spec *engine.Spec, step *engine.Step) []container.DeviceMapping {
	var to []container.DeviceMapping
	for _, mount := range step.Devices {
		device, ok := engine.LookupVolume(spec, mount.Name)
		if !ok {
			continue
		}
		if isDevice(device) == false {
			continue
		}
		to = append(to, container.DeviceMapping{
			PathOnHost:        device.HostPath.Path,
			PathInContainer:   mount.DevicePath,
			CgroupPermissions: "rwm",
		})
	}
	if len(to) == 0 {
		return nil
	}
	return to
}

// helper function returns a slice of volume mounts.
func toVolumeSlice(spec *engine.Spec, step *engine.Step) []string {
	// this entire function should be deprecated in
	// favor of toVolumeMounts, however, I am unable
	// to get it working with data volumes.
	var to []string
	for _, mount := range step.Volumes {
		volume, ok := engine.LookupVolume(spec, mount.Name)
		if !ok {
			continue
		}
		if isDevice(volume) {
			continue
		}
		if isDataVolume(volume) == false {
			continue
		}
		path := volume.Metadata.UID + ":" + mount.Path
		to = append(to, path)
	}
	return to
}

// helper function returns a slice of docker mount
// configurations.
func toVolumeMounts(spec *engine.Spec, step *engine.Step) []mount.Mount {
	var mounts []mount.Mount
	for _, target := range step.Volumes {
		source, ok := engine.LookupVolume(spec, target.Name)
		if !ok {
			continue
		}
		// HACK: this condition can be removed once
		// toVolumeSlice has been fully replaced. at this
		// time, I cannot figure out how to get mounts
		// working with data volumes :(
		if isDataVolume(source) {
			continue
		}
		mounts = append(mounts, toMount(source, target))
	}
	if len(mounts) == 0 {
		return nil
	}
	return mounts
}

// helper function converts the volume declaration to a
// docker mount structure.
func toMount(source *engine.Volume, target *engine.VolumeMount) mount.Mount {
	to := mount.Mount{
		Target: target.Path,
		Type:   toVolumeType(source),
	}
	if isBindMount(source) || isNamedPipe(source) {
		to.Source = source.HostPath.Path
	}
	if isTempfs(source) {
		to.TmpfsOptions = &mount.TmpfsOptions{
			SizeBytes: source.EmptyDir.SizeLimit,
			Mode:      0700,
		}
	}
	return to
}

// helper function returns the docker volume enumeration
// for the given volume.
func toVolumeType(from *engine.Volume) mount.Type {
	switch {
	case isDataVolume(from):
		return mount.TypeVolume
	case isTempfs(from):
		return mount.TypeTmpfs
	case isNamedPipe(from):
		return mount.TypeNamedPipe
	default:
		return mount.TypeBind
	}
}

// helper function that converts a key value map of
// environment variables to a string slice in key=value
// format.
func toEnv(env map[string]string) []string {
	var envs []string
	for k, v := range env {
		envs = append(envs, k+"="+v)
	}
	return envs
}

// returns true if the volume is a bind mount.
func isBindMount(volume *engine.Volume) bool {
	return volume.HostPath != nil
}

// returns true if the volume is in-memory.
func isTempfs(volume *engine.Volume) bool {
	return volume.EmptyDir != nil && volume.EmptyDir.Medium == "memory"
}

// returns true if the volume is a data-volume.
func isDataVolume(volume *engine.Volume) bool {
	return volume.EmptyDir != nil && volume.EmptyDir.Medium != "memory"
}

// returns true if the volume is a device
func isDevice(volume *engine.Volume) bool {
	return volume.HostPath != nil && strings.HasPrefix(volume.HostPath.Path, "/dev/")
}

// returns true if the volume is a named pipe.
func isNamedPipe(volume *engine.Volume) bool {
	return volume.HostPath != nil &&
		strings.HasPrefix(volume.HostPath.Path, `\\.\pipe\`)
}

// // helper function that converts a slice of device paths to a slice of
// // container.DeviceMapping.
// func toDevices(from []*engine.DeviceMapping) []container.DeviceMapping {
// 	var to []container.DeviceMapping
// 	for _, device := range from {
// 		to = append(to, container.DeviceMapping{
// 			PathOnHost:        device.Source,
// 			PathInContainer:   device.Target,
// 			CgroupPermissions: "rwm",
// 		})
// 	}
// 	return to
// }
