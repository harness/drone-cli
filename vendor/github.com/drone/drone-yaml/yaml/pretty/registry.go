package pretty

import (
	"github.com/drone/drone-yaml/yaml"
)

// helper function pretty prints the registry resource.
func printRegistry(w writer, v *yaml.Registry) {
	w.WriteString("---")
	w.WriteTagValue("kind", v.Kind)
	w.WriteTagValue("type", v.Type)
	if v.Type == "encrypted" {
		printData(w, v.Data)
	} else {
		w.WriteTagValue("data", v.Data)
	}
	w.WriteByte('\n')
	w.WriteByte('\n')
}
