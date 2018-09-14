package pretty

import (
	"github.com/drone/drone-yaml/yaml"
)

// helper function pretty prints the signature resource.
func printSignature(w writer, v *yaml.Signature) {
	w.WriteString("---")
	w.WriteTagValue("kind", v.Kind)
	w.WriteTagValue("hmac", v.Hmac)
	w.WriteByte('\n')
	w.WriteByte('\n')
}
