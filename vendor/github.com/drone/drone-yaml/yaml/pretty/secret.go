package pretty

import (
	"sort"
	"strings"

	"github.com/drone/drone-yaml/yaml"
)

// TODO consider "!!binary |" for secret value

// helper function to pretty prints the signature resource.
func printSecret(w writer, v *yaml.Secret) {
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

func printData(w writer, d map[string]string) {
	var keys []string
	for k := range d {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	w.WriteTag("data")
	w.IndentIncrease()
	for _, k := range keys {
		v := d[k]
		w.WriteTag(k)
		w.WriteByte(' ')
		w.WriteByte('>')
		w.IndentIncrease()
		v = spaceReplacer.Replace(v)
		for _, s := range chunk(v, 60) {
			w.WriteByte('\n')
			w.Indent()
			w.WriteString(s)
		}
		w.IndentDecrease()
	}
	w.IndentDecrease()
}

// replace spaces and newlines.
var spaceReplacer = strings.NewReplacer(" ", "", "\n", "")
