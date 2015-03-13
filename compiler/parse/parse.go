package parser

import (
	"github.com/drone/drone-cli/common"
	"github.com/drone/drone-cli/common/matrix"
	"gopkg.in/yaml.v2"
)

// Parse parses a build matrix and returns
// a list of build configurations for each axis.
func Parse(raw string) ([]*common.Config, error) {
	axis, err := matrix.Parse(raw)
	if err != nil {
		return nil, err
	}
	confs := []*common.Config{}

	// when no matrix values exist we should return
	// a single config value with an empty label.
	if len(axis) == 0 {
		conf, err := parse(raw)
		if err != nil {
			return nil, err
		}
		confs = append(confs, conf)
	}

	for _, ax := range axis {
		// inject the matrix values into the raw script
		injected := Inject(raw, ax)
		conf, err := parse(injected)
		if err != nil {
			return nil, err
		}
		conf.Axis = ax
		confs = append(confs, conf)
	}
	return confs, nil
}

// helper funtion to parse a yaml configuration file.
func parse(raw string) (*common.Config, error) {
	cfg := common.Config{}
	err := yaml.Unmarshal([]byte(raw), &cfg)
	return &cfg, err
}
