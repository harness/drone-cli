// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Drone Non-Commercial License
// that can be found in the LICENSE file.

package engine

import (
	"encoding/json"
	"io"
	"os"
	"strings"
)

// Parse parses the pipeline config from an io.Reader.
func Parse(r io.Reader) (*Spec, error) {
	cfg := Spec{}
	err := json.NewDecoder(r).Decode(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// ParseFile parses the pipeline config from a file.
func ParseFile(path string) (*Spec, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}

// ParseString parses the pipeline config from a string.
func ParseString(s string) (*Spec, error) {
	return Parse(
		strings.NewReader(s),
	)
}
