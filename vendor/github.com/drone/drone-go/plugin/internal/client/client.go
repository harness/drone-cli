// Copyright 2018 Drone.IO Inc.
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

package client

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/drone/drone-go/plugin/internal/aesgcm"

	"github.com/99designs/httpsignatures-go"
)

// required http headers
// note that (request-target) is disabled because reverse proxies,
// including aws lambda with api gateway, fail verification.
var headers = []string{
	"accept",
	"accept-encoding",
	"content-type",
	"date",
	"digest",
}

var signer = httpsignatures.NewSigner(
	httpsignatures.AlgorithmHmacSha256,
	headers...,
)

// New returns a new http.Client with signature verification.
func New(endpoint, secret string, skipverify bool) *Client {
	client := &Client{
		Accept:   "application/json",
		Encoding: "identity",
		Endpoint: endpoint,
		Secret:   secret,
	}
	if skipverify {
		client.Client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	}
	return client
}

// Client wraps an http.Client and applies retry logic and
// http signature verification.
type Client struct {
	Client     *http.Client
	Accept     string
	Encoding   string
	Endpoint   string
	Secret     string
	SkipVerify bool
}

// Do makes an http.Request to the target endpoint.
func (s *Client) Do(in, out interface{}) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	data, err := json.Marshal(in)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(data)
	req, err := http.NewRequest("POST", s.Endpoint, buf)
	if err != nil {
		return err
	}

	req = req.WithContext(ctx)
	req.Header.Add("Accept", s.Accept)
	req.Header.Add("Accept-Encoding", s.Encoding)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Digest", "SHA-256="+digest(data))
	req.Header.Add("Date", time.Now().UTC().Format(http.TimeFormat))
	err = signer.SignRequest("hmac-key", s.Secret, req)
	if err != nil {
		return err
	}

	res, err := s.client().Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode > 299 {
		// if the response body includes an error message
		// we should return the error string.
		if len(body) != 0 {
			return errors.New(
				string(body),
			)
		}
		// if the response body is empty we should return
		// the default status code text.
		return errors.New(
			http.StatusText(res.StatusCode),
		)
	}

	// if the response body return no content we exit
	// immediately. We do not read or unmarshal the response
	// and we do not return an error.
	if res.StatusCode == 204 {
		return nil
	}

	// the response body may be optionally encrypted
	// using the aesgcm algorithm. If encrypted,
	// decrypt using the shared secret.
	if res.Header.Get("Content-Encoding") == "aesgcm" {
		secret, err := aesgcm.Key(s.Secret)
		if err != nil {
			return err
		}
		plaintext, err := aesgcm.Decrypt(body, secret)
		if err != nil {
			return err
		}
		body = []byte(plaintext)
	}

	if out == nil {
		return nil
	}
	return json.Unmarshal(body, out)
}

func (s *Client) client() *http.Client {
	if s.Client == nil {
		return http.DefaultClient
	}
	return s.Client
}

func digest(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
