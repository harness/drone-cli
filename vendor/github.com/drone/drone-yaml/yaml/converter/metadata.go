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

package converter

// Metadata provides additional metadata used to
// convert the configuration file format.
type Metadata struct {
	// Filename of the configuration file, helps
	// determine the yaml configuration format.
	Filename string

	// Ref of the commit use to choose the correct
	// pipeline if the configuration format defines
	// multiple pipelines (like Bitbucket)
	Ref string
}
