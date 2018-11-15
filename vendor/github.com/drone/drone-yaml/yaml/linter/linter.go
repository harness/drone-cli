package linter

import (
	"errors"
	"fmt"

	"github.com/drone/drone-yaml/yaml"
)

// Lint performs lint operations for a resource.
func Lint(resource yaml.Resource, trusted bool) error {
	switch v := resource.(type) {
	case *yaml.Cron:
		return v.Validate()
	case *yaml.Pipeline:
		return checkPipeline(v, trusted)
	case *yaml.Secret:
		return v.Validate()
	case *yaml.Registry:
		return v.Validate()
	case *yaml.Signature:
		return v.Validate()
	default:
		return nil
	}
}

func checkPipeline(pipeline *yaml.Pipeline, trusted bool) error {
	err := checkVolumes(pipeline, trusted)
	if err != nil {
		return err
	}
	for _, container := range pipeline.Steps {
		err := checkContainer(container, trusted)
		if err != nil {
			return err
		}
	}
	for _, container := range pipeline.Services {
		err := checkContainer(container, trusted)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkContainer(container *yaml.Container, trusted bool) error {
	err := checkPorts(container.Ports, trusted)
	if err != nil {
		return err
	}
	if container.Build == nil && container.Image == "" {
		return errors.New("linter: invalid or missing image")
	}
	if container.Build != nil && container.Build.Image == "" {
		return errors.New("linter: invalid or missing build image")
	}
	if container.Name == "" {
		return errors.New("linter: invalid or missing name")
	}
	if trusted == false && container.Privileged {
		return errors.New("linter: untrusted repositories cannot enable privileged mode")
	}
	if trusted == false && len(container.Devices) > 0 {
		return errors.New("linter: untrusted repositories cannot mount devices")
	}
	if trusted == false && len(container.DNS) > 0 {
		return errors.New("linter: untrusted repositories cannot configure dns")
	}
	if trusted == false && len(container.DNSSearch) > 0 {
		return errors.New("linter: untrusted repositories cannot configure dns_search")
	}
	if trusted == false && len(container.ExtraHosts) > 0 {
		return errors.New("linter: untrusted repositories cannot configure extra_hosts")
	}
	for _, mount := range container.Volumes {
		switch mount.Name {
		case "workspace", "_workspace", "_docker_socket":
			return fmt.Errorf("linter: invalid volume name: %s", mount.Name)
		}
	}
	return nil
}

func checkPorts(ports []*yaml.Port, trusted bool) error {
	for _, port := range ports {
		err := checkPort(port, trusted)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkPort(port *yaml.Port, trusted bool) error {
	if trusted == false && port.Host != 0 {
		return errors.New("linter: untrusted repositories cannot map to a host port")
	}
	return nil
}

func checkVolumes(pipeline *yaml.Pipeline, trusted bool) error {
	for _, volume := range pipeline.Volumes {
		if volume.EmptyDir != nil {
			err := checkEmptyDirVolume(volume.EmptyDir, trusted)
			if err != nil {
				return err
			}
		}
		if volume.HostPath != nil {
			err := checkHostPathVolume(volume.HostPath, trusted)
			if err != nil {
				return err
			}
		}
		switch volume.Name {
		case "workspace", "_workspace", "_docker_socket":
			return fmt.Errorf("linter: invalid volume name: %s", volume.Name)
		}
	}
	return nil
}

func checkHostPathVolume(volume *yaml.VolumeHostPath, trusted bool) error {
	if trusted == false {
		return errors.New("linter: untrusted repsitories cannot mount host volumes")
	}
	return nil
}

func checkEmptyDirVolume(volume *yaml.VolumeEmptyDir, trusted bool) error {
	if trusted == false && volume.Medium == "memory" {
		return errors.New("linter: untrusted repsitories cannot mount in-memory volumes")
	}
	return nil
}
