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
