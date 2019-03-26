// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Drone Non-Commercial License
// that can be found in the LICENSE file.

// +build !oss

package converter

import (
	"github.com/drone/drone-yaml/yaml/converter/bitbucket"
	"github.com/drone/drone-yaml/yaml/converter/gitlab"
	"github.com/drone/drone-yaml/yaml/converter/legacy"
)

// Convert converts the yaml configuration file from
// the legacy format to the 1.0+ format.
func Convert(d []byte, m Metadata) ([]byte, error) {
	switch m.Filename {
	case "bitbucket-pipelines.yml":
		return bitbucket.Convert(d, m.Ref)
	case "circle.yml", ".circleci/config.yml":
		// TODO(bradrydzewski)
	case ".gitlab-ci.yml":
		return gitlab.Convert(d)
	case ".travis.yml":
		// TODO(bradrydzewski)
	}
	// if the filename does not match any external
	// systems we check to see if the configuration
	// file is a legacy (pre 1.0) .drone.yml format.
	if legacy.Match(d) {
		return legacy.Convert(d)
	}
	// else return the unmodified configuration
	// back to the caller.
	return d, nil
}

// ConvertString converts the yaml configuration file from
// the legacy format to the 1.0+ format.
func ConvertString(s string, m Metadata) (string, error) {
	b, err := Convert([]byte(s), m)
	return string(b), err
}
