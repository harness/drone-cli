// Copyright 2019 Drone IO, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/drone/drone-runtime/engine"
)

// config represents the Docker client configuration,
// typically located at ~/.docker/config.json
type config struct {
	Auths map[string]auths `json:"auths"`
}

type auths struct {
	Auth string `json:"auth"`
}

// Parse parses the registry credential from the reader.
func Parse(r io.Reader) ([]*engine.DockerAuth, error) {
	c := new(config)
	err := json.NewDecoder(r).Decode(c)
	if err != nil {
		return nil, err
	}
	var auths []*engine.DockerAuth
	for k, v := range c.Auths {
		username, password := decode(v.Auth)
		auths = append(auths, &engine.DockerAuth{
			Address:  hostname(k),
			Username: username,
			Password: password,
		})
	}
	return auths, nil
}

// ParseFile parses the registry credential file.
func ParseFile(filepath string) ([]*engine.DockerAuth, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}

// ParseString parses the registry credential file.
func ParseString(s string) ([]*engine.DockerAuth, error) {
	return Parse(strings.NewReader(s))
}

// encode returns the encoded credentials.
func encode(username, password string) string {
	return base64.StdEncoding.EncodeToString(
		[]byte(username + ":" + password),
	)
}

// decode returns the decoded credentials.
func decode(s string) (username, password string) {
	d, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return
	}
	parts := strings.SplitN(string(d), ":", 2)
	if len(parts) > 0 {
		username = parts[0]
	}
	if len(parts) > 1 {
		password = parts[1]
	}
	return
}

func hostname(s string) string {
	uri, _ := url.Parse(s)
	if uri.Host != "" {
		s = uri.Host
	}
	return s
}

// Encode returns the json marshaled, base64 encoded
// credential string that can be passed to the docker
// registry authentication header.
func Encode(username, password string) string {
	v := struct {
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}{
		Username: username,
		Password: password,
	}
	buf, _ := json.Marshal(&v)
	return base64.URLEncoding.EncodeToString(buf)
}

// Marshal marshals the DockerAuth credentials to a
// .docker/config.json file.
func Marshal(list []*engine.DockerAuth) ([]byte, error) {
	out := &config{}
	out.Auths = map[string]auths{}
	for _, item := range list {
		out.Auths[item.Address] = auths{
			Auth: encode(
				item.Username,
				item.Password,
			),
		}
	}
	return json.Marshal(out)
}
