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
