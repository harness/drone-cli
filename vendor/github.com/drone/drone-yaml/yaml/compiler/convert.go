// Copyright the Drone Authors.
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
)

func toIgnoreErr(from *yaml.Container) bool {
	return strings.EqualFold(from.Failure, "ignore")
}

func toPorts(from *yaml.Container) []*engine.Port {
	var ports []*engine.Port
	for _, port := range from.Ports {
		ports = append(ports, toPort(port))
	}
	return ports
}

func toPort(from *yaml.Port) *engine.Port {
	return &engine.Port{
		Port:     from.Port,
		Host:     from.Host,
		Protocol: from.Protocol,
	}
}

func toPullPolicy(from *yaml.Container) engine.PullPolicy {
	switch strings.ToLower(from.Pull) {
	case "always":
		return engine.PullAlways
	case "if-not-exists":
		return engine.PullIfNotExists
	case "never":
		return engine.PullNever
	default:
		return engine.PullDefault
	}
}

func toRunPolicy(from *yaml.Container) engine.RunPolicy {
	onFailure := from.When.Status.Match("failure") && len(from.When.Status.Include) > 0
	onSuccess := from.When.Status.Match("success")
	switch {
	case onFailure && onSuccess:
		return engine.RunAlways
	case onFailure:
		return engine.RunOnFailure
	case onSuccess:
		return engine.RunOnSuccess
	default:
		return engine.RunNever
	}
}

func toResources(from *yaml.Container) *engine.Resources {
	if from.Resources == nil {
		return nil
	}
	return &engine.Resources{
		Limits:   toResourceObject(from.Resources.Limits),
		Requests: toResourceObject(from.Resources.Requests),
	}
}

func toResourceObject(from *yaml.ResourceObject) *engine.ResourceObject {
	if from == nil {
		return nil
	}
	return &engine.ResourceObject{
		CPU:    toCPUMillis(from.CPU),
		Memory: int64(from.Memory),
	}
}

func toCPUMillis(f float64) int64 {
	if f > 0 {
		f *= 1000
	}
	return int64(f)
}
