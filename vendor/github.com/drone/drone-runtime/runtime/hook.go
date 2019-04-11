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

package runtime

// Hook provides a set of hooks to run at various stages of
// runtime execution.
type Hook struct {
	// Before is called before all all steps are executed.
	Before func(*State) error

	// BeforeEach is called before each step is executed.
	BeforeEach func(*State) error

	// After is called after all steps are executed.
	After func(*State) error

	// AfterEach is called after each step is executed.
	AfterEach func(*State) error

	// GotLine is called when a line is logged.
	GotLine func(*State, *Line) error

	// GotLogs is called when the logs are completed.
	GotLogs func(*State, []*Line) error
}
