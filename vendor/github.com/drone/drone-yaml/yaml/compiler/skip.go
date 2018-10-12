package compiler

import "github.com/drone/drone-yaml/yaml"

// SkipData provides build metadata use to determine if a
// pipeline step should be skipped.
type SkipData struct {
	Branch   string
	Event    string
	Instance string
	Ref      string
	Repo     string
	Target   string
}

// SkipFunc returns a function that can be used to skip
// individual pipeline steps based on build metadata.
func SkipFunc(data SkipData) func(*yaml.Container) bool {
	return func(container *yaml.Container) bool {
		switch {
		case !container.When.Branch.Match(data.Branch):
			return true
		case !container.When.Event.Match(data.Event):
			return true
		case !container.When.Instance.Match(data.Instance):
			return true
		case !container.When.Ref.Match(data.Ref):
			return true
		case !container.When.Repo.Match(data.Repo):
			return true
		case !container.When.Target.Match(data.Target):
			return true
		default:
			return false
		}
	}
}
