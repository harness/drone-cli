package converter

import (
	"github.com/drone/drone-yaml/yaml/converter/internal"
)

// ConvertBytes converts the yaml configuration file from
// the legacy format to the 1.0+ format.
func ConvertBytes(d []byte) ([]byte, error) {
	return yaml.ConvertBytes(d)
}

// ConvertString converts the yaml configuration file from
// the legacy format to the 1.0+ format.
func ConvertString(s string) (string, error) {
	return yaml.ConvertString(s)
}
