package config

import (
	"sort"
	"strings"
)

func Inject(config string, params map[string]string) string {
	if params == nil {
		return config
	}
	keys := []string{}
	for k, _ := range params {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	injected := config
	for _, k := range keys {
		v := params[k]
		injected = strings.Replace(injected, "$$"+k, v, -1)
	}
	return injected
}
