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

package transform

import (
	"os"
	"strings"

	"github.com/drone/drone-runtime/engine"
)

// WithProxy is a transform function that adds the
// http_proxy environment variables to every container.
func WithProxy() func(*engine.Spec) {
	environ := map[string]string{}
	if value := getenv("no_proxy"); value != "" {
		environ["no_proxy"] = value
		environ["NO_PROXY"] = value
	}
	if value := getenv("http_proxy"); value != "" {
		environ["http_proxy"] = value
		environ["HTTP_PROXY"] = value
	}
	if value := getenv("https_proxy"); value != "" {
		environ["https_proxy"] = value
		environ["HTTPS_PROXY"] = value
	}
	return WithEnviron(environ)
}

func getenv(name string) (value string) {
	name = strings.ToUpper(name)
	if value := os.Getenv(name); value != "" {
		return value
	}
	name = strings.ToLower(name)
	if value := os.Getenv(name); value != "" {
		return value
	}
	return
}
