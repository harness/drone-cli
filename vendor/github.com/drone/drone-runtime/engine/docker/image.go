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

	"github.com/docker/distribution/reference"
)

// helper function parses the image and returns the
// canonical image name, domain name, and whether or not
// the image tag is :latest.
func parseImage(s string) (canonical, domain string, latest bool, err error) {
	// parse the docker image name. We need to extract the
	// image domain name and match to registry credentials
	// stored in the .docker/config.json object.
	named, err := reference.ParseNormalizedNamed(s)
	if err != nil {
		return
	}
	// the canonical image name, for some reason, excludes
	// the tag name. So we need to make sure it is included
	// in the image name so we can determine if the :latest
	// tag is specified
	named = reference.TagNameOnly(named)

	return named.String(),
		reference.Domain(named),
		strings.HasSuffix(named.String(), ":latest"),
		nil
}
