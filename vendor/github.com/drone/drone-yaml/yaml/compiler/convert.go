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
		CPU:    int64(from.CPU),
		Memory: int64(from.Memory),
	}
}
