// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Drone Non-Commercial License
// that can be found in the LICENSE file.

package legacy

import "github.com/drone/drone-yaml/yaml/converter/legacy/internal"

// Convert converts the yaml configuration file from
// the legacy format to the 1.0+ format.
func Convert(d []byte, remote string) ([]byte, error) {
	return yaml.Convert(d, remote)
}
