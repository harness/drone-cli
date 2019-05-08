// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Drone Non-Commercial License
// that can be found in the LICENSE file.

package yaml

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	droneyaml "github.com/drone/drone-yaml/yaml"
	"github.com/drone/drone-yaml/yaml/converter/legacy/matrix"
	"github.com/drone/drone-yaml/yaml/pretty"

	"gopkg.in/yaml.v2"
)

// Config provides the high-level configuration.
type Config struct {
	Workspace struct {
		Base string
		Path string
	}
	Clone    Containers
	Pipeline Containers
	Services Containers
	Branches Constraint
	Matrix   interface{}
	Secrets  map[string]struct {
		Driver     string
		DriverOpts map[string]string `yaml:"driver_opts"`
		Path       string
		Vault      string
	}
}

// Convert converts the yaml configuration file from
// the legacy format to the 1.0+ format.
func Convert(d []byte) ([]byte, error) {
	from := new(Config)
	err := yaml.Unmarshal(d, from)

	if err != nil {
		return nil, err
	}

	manifest := &droneyaml.Manifest{}

	pipeline := droneyaml.Pipeline{}
	pipeline.Name = "default"
	pipeline.Kind = "pipeline"

	pipeline.Workspace.Base = from.Workspace.Base
	pipeline.Workspace.Path = from.Workspace.Path

	if len(from.Clone.Containers) != 0 {
		pipeline.Clone.Disable = true
		for _, container := range from.Clone.Containers {
			pipeline.Steps = append(pipeline.Steps,
				toContainer(container),
			)
		}
	}

	for _, container := range from.Services.Containers {
		pipeline.Services = append(pipeline.Services,
			toContainer(container),
		)
	}

	for _, container := range from.Pipeline.Containers {
		pipeline.Steps = append(pipeline.Steps,
			toContainer(container),
		)
	}

	pipeline.Volumes = toVolumes(from)
	pipeline.Trigger.Branch.Include = from.Branches.Include
	pipeline.Trigger.Branch.Exclude = from.Branches.Exclude

	if from.Matrix != nil {
		axes, err := matrix.Parse(d)

		if err != nil {
			return nil, err
		}

		for index, environ := range axes {
			current := pipeline
			current.Name = fmt.Sprintf("matrix-%d", index+1)

			services := make([]*droneyaml.Container, 0)
			for _, service := range current.Services {
				if len(service.When.Matrix) == 0 {
					services = append(services, service)
					continue
				}

				for whenKey, whenValue := range service.When.Matrix {
					for envKey, envValue := range environ {
						if whenKey == envKey && whenValue == envValue {
							services = append(services, service)
						}
					}
				}
			}
			current.Services = services

			steps := make([]*droneyaml.Container, 0)
			for _, step := range current.Steps {
				if len(step.When.Matrix) == 0 {
					steps = append(steps, step)
					continue
				}

				for whenKey, whenValue := range step.When.Matrix {
					for envKey, envValue := range environ {
						if whenKey == envKey && whenValue == envValue {
							steps = append(steps, step)
						}
					}
				}
			}
			current.Steps = steps

			marshaled, err := yaml.Marshal(&current)

			if err != nil {
				return nil, err
			}

			transformed := string(marshaled)

			for key, value := range environ {
				if strings.Contains(value, "\n") {
					value = fmt.Sprintf("%q", value)
				}

				transformed = strings.Replace(transformed, fmt.Sprintf("${%s}", key), value, -1)
			}

			if err := yaml.Unmarshal([]byte(transformed), &current); err != nil {
				return nil, err
			}

			manifest.Resources = append(manifest.Resources, &current)
		}
	} else {
		manifest.Resources = append(manifest.Resources, &pipeline)
	}

	secrets := toSecrets(from)
	for _, secret := range secrets {
		manifest.Resources = append(manifest.Resources, secret)
	}

	buf := new(bytes.Buffer)
	pretty.Print(buf, manifest)

	return buf.Bytes(), nil
}

func toContainer(from *Container) *droneyaml.Container {
	return &droneyaml.Container{
		Name:        from.Name,
		Image:       from.Image,
		Detach:      from.Detached,
		Command:     from.Command,
		Commands:    from.Commands,
		DNS:         from.DNS,
		DNSSearch:   from.DNSSearch,
		Entrypoint:  from.Entrypoint,
		Environment: toEnvironment(from),
		ExtraHosts:  from.ExtraHosts,
		Pull:        toPullPolicy(from.Pull),
		Privileged:  from.Privileged,
		Settings:    toSettings(from.Vargs),
		Volumes:     toVolumeMounts(from.Volumes),
		When:        toConditions(from.Constraints),
	}
}

// helper function converts the legacy constraint syntax
// to the new condition syntax.
func toConditions(from Constraints) droneyaml.Conditions {
	return droneyaml.Conditions{
		Ref: droneyaml.Condition{
			Include: from.Ref.Include,
			Exclude: from.Ref.Exclude,
		},
		Repo: droneyaml.Condition{
			Include: from.Repo.Include,
			Exclude: from.Repo.Exclude,
		},
		Instance: droneyaml.Condition{
			Include: from.Instance.Include,
			Exclude: from.Instance.Exclude,
		},
		Target: droneyaml.Condition{
			Include: from.Environment.Include,
			Exclude: from.Environment.Exclude,
		},
		Event: droneyaml.Condition{
			Include: from.Event.Include,
			Exclude: from.Event.Exclude,
		},
		Branch: droneyaml.Condition{
			Include: from.Branch.Include,
			Exclude: from.Branch.Exclude,
		},
		Status: droneyaml.Condition{
			Include: from.Status.Include,
			Exclude: from.Status.Exclude,
		},
		Matrix: from.Matrix,
	}
}

// helper function converts the legacy environment syntax
// to the new environment syntax.
func toEnvironment(from *Container) map[string]*droneyaml.Variable {
	envs := map[string]*droneyaml.Variable{}
	for key, val := range from.Environment.Map {
		envs[key] = &droneyaml.Variable{
			Value: val,
		}
	}
	for _, val := range from.Secrets.Secrets {
		name := strings.ToUpper(val.Target)
		envs[name] = &droneyaml.Variable{
			Secret: val.Source,
		}
	}
	return envs
}

// helper function converts the legacy image pull syntax
// to the new pull policy syntax.
func toPullPolicy(pull bool) string {
	switch pull {
	case true:
		return "always"
	default:
		return "default"
	}
}

// helper function converts the legacy secret syntax to the
// new secret variable syntax.
func toSecrets(from *Config) []*droneyaml.Secret {
	var keys []string
	for key := range from.Secrets {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var secrets []*droneyaml.Secret
	for _, key := range keys {
		val := from.Secrets[key]
		secret := new(droneyaml.Secret)
		secret.Name = key
		secret.Kind = "secret"

		if val.Driver == "vault" {
			if val.DriverOpts != nil {
				secret.Get.Path = val.DriverOpts["path"]
				secret.Get.Name = val.DriverOpts["key"]
			}
		} else if val.Path != "" {
			secret.Get.Path = val.Path
		} else {
			secret.Get.Path = val.Vault
		}
		secrets = append(secrets, secret)
	}
	if len(secrets) == 0 {
		return nil
	}
	return secrets
}

// helper function converts the legacy vargs syntax to the
// new environment syntax.
func toSettings(from map[string]interface{}) map[string]*droneyaml.Parameter {
	params := map[string]*droneyaml.Parameter{}
	for key, val := range from {
		params[key] = &droneyaml.Parameter{
			Value: val,
		}
	}
	return params
}

// helper function converts the legacy volume syntax
// to the new volume mount syntax.
func toVolumeMounts(from []*Volume) []*droneyaml.VolumeMount {
	to := []*droneyaml.VolumeMount{}
	for _, v := range from {
		to = append(to, &droneyaml.VolumeMount{
			Name:      fmt.Sprintf("%x", v.Source),
			MountPath: v.Destination,
		})
	}
	return to
}

// helper function converts the legacy volume syntax
// to the new volume mount syntax.
func toVolumes(from *Config) []*droneyaml.Volume {
	set := map[string]struct{}{}
	to := []*droneyaml.Volume{}

	containers := []*Container{}
	containers = append(containers, from.Pipeline.Containers...)
	containers = append(containers, from.Services.Containers...)

	for _, container := range containers {
		for _, v := range container.Volumes {
			name := fmt.Sprintf("%x", v.Source)
			if _, ok := set[name]; ok {
				continue
			}
			set[name] = struct{}{}
			to = append(to, &droneyaml.Volume{
				Name: name,
				HostPath: &droneyaml.VolumeHostPath{
					Path: v.Source,
				},
			})
		}
	}
	return to
}
