package config

import (
	"github.com/drone/drone-cli/common"
	"github.com/drone/drone-cli/common/matrix"
	"gopkg.in/yaml.v2"
)

// Parse parses a yaml configuration file.
func Parse(raw string) (*common.Config, error) {
	cfg := common.Config{}
	err := yaml.Unmarshal([]byte(raw), &cfg)
	return &cfg, err
}

// ParseMatrix parses a build matrix and returns
// a list of build configurations for each axis.
func ParseMatrix(raw string) (map[string]*common.Config, error) {
	axis, err := matrix.Parse(raw)
	if err != nil {
		return nil, err
	}
	confs := map[string]*common.Config{}

	// when no matrix values exist we should return
	// a single config value with an empty label.
	if len(axis) == 0 {
		conf, err := Parse(raw)
		confs[""] = conf
		return confs, err
	}

	for _, a := range axis {
		// inject the matrix values into the raw script
		injected := Inject(raw, a)
		cfg, err := Parse(injected)
		if err != nil {
			return nil, err
		}
		// append the matrix values to the env variables
		// for k, v := range a {
		// 	cfg.Environment[k] = v
		// }
		confs[a.String()] = cfg
	}
	return confs, nil
}
