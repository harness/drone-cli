package legacy

import "github.com/drone/drone-yaml/yaml/converter/legacy/internal"

// Convert converts the yaml configuration file from
// the legacy format to the 1.0+ format.
func Convert(d []byte) ([]byte, error) {
	return yaml.Convert(d)
}
