package converter

import (
	"regexp"
)

var re = regexp.MustCompile(`(?m)^pipeline:(\s+)?$`)

// IsLegacy returns true if the yaml configuration file is
// legacy and requires converstion.
func IsLegacy(s string) bool {
	matches := re.FindAllString(s, -1)
	return len(matches) != 0
}

// IsLegacyBytes returns true if the yaml configuration file
// is legacy and requires converstion.
func IsLegacyBytes(b []byte) bool {
	matches := re.FindAll(b, -1)
	return len(matches) != 0
}
