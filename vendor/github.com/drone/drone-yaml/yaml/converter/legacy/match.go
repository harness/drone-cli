// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Drone Non-Commercial License
// that can be found in the LICENSE file.

package legacy

import (
	"regexp"
)

var re = regexp.MustCompile(`(?m)^pipeline:(\s+)?$`)

// Match returns true if the yaml configuration file
// is legacy and requires converstion.
func Match(b []byte) bool {
	matches := re.FindAll(b, -1)
	return len(matches) != 0
}
