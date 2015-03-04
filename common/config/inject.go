package config

import (
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

func Inject(raw string, params map[string]string) string {
	if params == nil {
		return raw
	}
	keys := []string{}
	for k, _ := range params {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	injected := raw
	for _, k := range keys {
		v := params[k]
		injected = strings.Replace(injected, "$$"+k, v, -1)
	}
	return injected
}

func InjectSafe(raw string, params map[string]string) string {
	before, _ := Parse(raw)
	after, _ := Parse(Inject(raw, params))
	after.Build = before.Build
	after.Compose = before.Compose
	scrubbed, _ := yaml.Marshal(after)
	return string(scrubbed)
}
