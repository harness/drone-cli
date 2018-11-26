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

package drone

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	pathSelf            = "%s/api/user"
	pathFeed            = "%s/api/user/feed"
	pathRepos           = "%s/api/user/repos"
	pathRepo            = "%s/api/repos/%s/%s"
	pathRepoMove        = "%s/api/repos/%s/%s/move?to=%s"
	pathChown           = "%s/api/repos/%s/%s/chown"
	pathRepair          = "%s/api/repos/%s/%s/repair"
	pathBuilds          = "%s/api/repos/%s/%s/builds?%s"
	pathBuild           = "%s/api/repos/%s/%s/builds/%v"
	pathApprove         = "%s/api/repos/%s/%s/builds/%d/approve/%d"
	pathDecline         = "%s/api/repos/%s/%s/builds/%d/decline/%d"
	pathPromote         = "%s/api/repos/%s/%s/builds/%d/promote?%s"
	pathRollback        = "%s/api/repos/%s/%s/builds/%d/rollback?%s"
	pathJob             = "%s/api/repos/%s/%s/builds/%d/%d"
	pathLog             = "%s/api/repos/%s/%s/builds/%d/logs/%d/%d"
	pathRepoSecrets     = "%s/api/repos/%s/%s/secrets"
	pathRepoSecret      = "%s/api/repos/%s/%s/secrets/%s"
	pathRepoRegistries  = "%s/api/repos/%s/%s/registry"
	pathRepoRegistry    = "%s/api/repos/%s/%s/registry/%s"
	pathEncryptSecret   = "%s/api/repos/%s/%s/encrypt/secret"
	pathEncryptRegistry = "%s/api/repos/%s/%s/encrypt/registry"
	pathSign            = "%s/api/repos/%s/%s/sign"
	pathVerify          = "%s/api/repos/%s/%s/verify"
	pathCrons           = "%s/api/repos/%s/%s/cron"
	pathCron            = "%s/api/repos/%s/%s/cron/%s"
	pathUsers           = "%s/api/users"
	pathUser            = "%s/api/users/%s"
	pathQueue           = "%s/api/queue"
	pathServers         = "%s/api/servers"
	pathServer          = "%s/api/servers/%s"
	pathScalerPause     = "%s/api/pause"
	pathScalerResume    = "%s/api/resume"
	pathNodes           = "%s/api/nodes"
	pathNode            = "%s/api/nodes/%s"
	pathVersion         = "%s/version"
)

type client struct {
	client *http.Client
	addr   string
}

type ListOptions struct {
	Page int
}

func encodeListOptions(opts ListOptions) string {
	params := url.Values{}
	if opts.Page != 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	return params.Encode()
}

// New returns a client at the specified url.
func New(uri string) Client {
	return &client{http.DefaultClient, strings.TrimSuffix(uri, "/")}
}

// NewClient returns a client at the specified url.
func NewClient(uri string, cli *http.Client) Client {
	return &client{cli, strings.TrimSuffix(uri, "/")}
}

// SetClient sets the http.Client.
func (c *client) SetClient(client *http.Client) {
	c.client = client
}

// SetAddress sets the server address.
func (c *client) SetAddress(addr string) {
	c.addr = addr
}

// Self returns the currently authenticated user.
func (c *client) Self() (*User, error) {
	out := new(User)
	uri := fmt.Sprintf(pathSelf, c.addr)
	err := c.get(uri, out)
	return out, err
}

// User returns a user by login.
func (c *client) User(login string) (*User, error) {
	out := new(User)
	uri := fmt.Sprintf(pathUser, c.addr, login)
	err := c.get(uri, out)
	return out, err
}

// UserList returns a list of all registered users.
func (c *client) UserList() ([]*User, error) {
	var out []*User
	uri := fmt.Sprintf(pathUsers, c.addr)
	err := c.get(uri, &out)
	return out, err
}

// UserCreate creates a new user account.
func (c *client) UserCreate(in *User) (*User, error) {
	out := new(User)
	uri := fmt.Sprintf(pathUsers, c.addr)
	err := c.post(uri, in, out)
	return out, err
}

// UserUpdate updates a user account.
func (c *client) UserUpdate(in *User) (*User, error) {
	out := new(User)
	uri := fmt.Sprintf(pathUser, c.addr, in.Login)
	err := c.patch(uri, in, out)
	return out, err
}

// UserDelete deletes a user account.
func (c *client) UserDelete(login string) error {
	uri := fmt.Sprintf(pathUser, c.addr, login)
	err := c.delete(uri)
	return err
}

// Repo returns a repository by name.
func (c *client) Repo(owner string, name string) (*Repo, error) {
	out := new(Repo)
	uri := fmt.Sprintf(pathRepo, c.addr, owner, name)
	err := c.get(uri, out)
	return out, err
}

// RepoList returns a list of all repositories to which
// the user has explicit access in the host system.
func (c *client) RepoList() ([]*Repo, error) {
	var out []*Repo
	uri := fmt.Sprintf(pathRepos, c.addr)
	err := c.get(uri, &out)
	return out, err
}

// RepoListSync returns a list of all repositories to which
// the user has explicit access in the host system.
func (c *client) RepoListSync() ([]*Repo, error) {
	var out []*Repo
	uri := fmt.Sprintf(pathRepos, c.addr)
	err := c.post(uri, nil, &out)
	return out, err
}

// RepoEnable activates a repository.
func (c *client) RepoEnable(owner, name string) (*Repo, error) {
	out := new(Repo)
	uri := fmt.Sprintf(pathRepo, c.addr, owner, name)
	err := c.post(uri, nil, out)
	return out, err
}

// RepoDisable disables a repository.
func (c *client) RepoDisable(owner, name string) error {
	uri := fmt.Sprintf(pathRepo, c.addr, owner, name)
	err := c.delete(uri)
	return err
}

// RepoUpdate updates a repository.
func (c *client) RepoUpdate(owner, name string, in *RepoPatch) (*Repo, error) {
	out := new(Repo)
	uri := fmt.Sprintf(pathRepo, c.addr, owner, name)
	err := c.patch(uri, in, out)
	return out, err
}

// RepoChown updates a repository owner.
func (c *client) RepoChown(owner, name string) (*Repo, error) {
	out := new(Repo)
	uri := fmt.Sprintf(pathChown, c.addr, owner, name)
	err := c.post(uri, nil, out)
	return out, err
}

// RepoRepair repais the repository hooks.
func (c *client) RepoRepair(owner, name string) error {
	uri := fmt.Sprintf(pathRepair, c.addr, owner, name)
	return c.post(uri, nil, nil)
}

// Build returns a repository build by number.
func (c *client) Build(owner, name string, num int) (*Build, error) {
	out := new(Build)
	uri := fmt.Sprintf(pathBuild, c.addr, owner, name, num)
	err := c.get(uri, out)
	return out, err
}

// Build returns the latest repository build by branch.
func (c *client) BuildLast(owner, name, branch string) (*Build, error) {
	out := new(Build)
	uri := fmt.Sprintf(pathBuild, c.addr, owner, name, "latest")
	if len(branch) != 0 {
		uri += "?branch=" + branch
	}
	err := c.get(uri, out)
	return out, err
}

// BuildList returns a list of recent builds for the
// the specified repository.
func (c *client) BuildList(owner, name string, opts ListOptions) ([]*Build, error) {
	var out []*Build
	uri := fmt.Sprintf(pathBuilds, c.addr, owner, name, encodeListOptions(opts))
	err := c.get(uri, &out)
	return out, err
}

// BuildRestart re-starts a stopped build.
func (c *client) BuildRestart(owner, name string, build int, params map[string]string) (*Build, error) {
	out := new(Build)
	val := mapValues(params)
	uri := fmt.Sprintf(pathBuild, c.addr, owner, name, build)
	if len(params) > 0 {
		uri = uri + "?" + val.Encode()
	}
	err := c.post(uri, nil, out)
	return out, err
}

// BuildCancel cancels the running job.
func (c *client) BuildCancel(owner, name string, build int) error {
	uri := fmt.Sprintf(pathBuild, c.addr, owner, name, build)
	err := c.delete(uri)
	return err
}

// Promote promotes a build to the target environment.
func (c *client) Promote(namespace, name string, build int, target string, params map[string]string) (*Build, error) {
	out := new(Build)
	val := mapValues(params)
	val.Set("target", target)
	uri := fmt.Sprintf(pathPromote, c.addr, namespace, name, build, val.Encode())
	err := c.post(uri, nil, out)
	return out, err
}

// Roolback reverts the target environment to an previous build.
func (c *client) Rollback(namespace, name string, build int, target string, params map[string]string) (*Build, error) {
	out := new(Build)
	val := mapValues(params)
	val.Set("target", target)
	uri := fmt.Sprintf(pathRollback, c.addr, namespace, name, build, val.Encode())
	err := c.post(uri, nil, out)
	return out, err
}

// Approve approves a blocked build stage.
func (c *client) Approve(namespace, name string, build, stage int) error {
	uri := fmt.Sprintf(pathApprove, c.addr, namespace, name, build, stage)
	err := c.post(uri, nil, nil)
	return err
}

// Decline declines a blocked build stage.
func (c *client) Decline(namespace, name string, build, stage int) error {
	uri := fmt.Sprintf(pathDecline, c.addr, namespace, name, build, stage)
	err := c.post(uri, nil, nil)
	return err
}

// BuildLogs returns the build logs for the specified job.
func (c *client) Logs(owner, name string, build, stage, step int) ([]*Line, error) {
	var out []*Line
	uri := fmt.Sprintf(pathLog, c.addr, owner, name, build, stage, step)
	err := c.get(uri, &out)
	return out, err
}

// LogsPurge purges the build logs for the specified build.
func (c *client) LogsPurge(owner, name string, build, stage, step int) error {
	uri := fmt.Sprintf(pathLog, c.addr, owner, name, build, stage, step)
	err := c.delete(uri)
	return err
}

// Sign signs the yaml file.
func (c *client) Sign(owner, name, file string) (string, error) {
	in := struct {
		Data string `json:"data"`
	}{Data: file}
	out := struct {
		Data string `json:"data"`
	}{}
	uri := fmt.Sprintf(pathSign, c.addr, owner, name)
	err := c.post(uri, &in, &out)
	return out.Data, err
}

// Verify verifies the yaml signature.
func (c *client) Verify(owner, name, file string) error {
	in := struct {
		Data string `json:"data"`
	}{Data: file}
	uri := fmt.Sprintf(pathVerify, c.addr, owner, name)
	return c.post(uri, &in, nil)
}

// Encrypt returns an encrypted secret.
func (c *client) Encrypt(owner, name string, secret *Secret) (string, error) {
	out := struct {
		Data string `json:"data"`
	}{}
	uri := fmt.Sprintf(pathEncryptSecret, c.addr, owner, name)
	err := c.post(uri, secret, &out)
	return out.Data, err
}

// Secret returns a secret by name.
func (c *client) Secret(owner, name, secret string) (*Secret, error) {
	out := new(Secret)
	uri := fmt.Sprintf(pathRepoSecret, c.addr, owner, name, secret)
	err := c.get(uri, out)
	return out, err
}

// SecretList returns a list of all repository secrets.
func (c *client) SecretList(owner string, name string) ([]*Secret, error) {
	var out []*Secret
	uri := fmt.Sprintf(pathRepoSecrets, c.addr, owner, name)
	err := c.get(uri, &out)
	return out, err
}

// SecretCreate creates a secret.
func (c *client) SecretCreate(owner, name string, in *Secret) (*Secret, error) {
	out := new(Secret)
	uri := fmt.Sprintf(pathRepoSecrets, c.addr, owner, name)
	err := c.post(uri, in, out)
	return out, err
}

// SecretUpdate updates a secret.
func (c *client) SecretUpdate(owner, name string, in *Secret) (*Secret, error) {
	out := new(Secret)
	uri := fmt.Sprintf(pathRepoSecret, c.addr, owner, name, in.Name)
	err := c.patch(uri, in, out)
	return out, err
}

// SecretDelete deletes a secret.
func (c *client) SecretDelete(owner, name, secret string) error {
	uri := fmt.Sprintf(pathRepoSecret, c.addr, owner, name, secret)
	return c.delete(uri)
}

// Cron returns a cronjob by name.
func (c *client) Cron(owner, name, cron string) (*Cron, error) {
	out := new(Cron)
	uri := fmt.Sprintf(pathCron, c.addr, owner, name, cron)
	err := c.get(uri, out)
	return out, err
}

// CronList returns a list of all repository cronjobs.
func (c *client) CronList(owner string, name string) ([]*Cron, error) {
	var out []*Cron
	uri := fmt.Sprintf(pathCrons, c.addr, owner, name)
	err := c.get(uri, &out)
	return out, err
}

// CronCreate creates a cronjob.
func (c *client) CronCreate(owner, name string, in *Cron) (*Cron, error) {
	out := new(Cron)
	uri := fmt.Sprintf(pathCrons, c.addr, owner, name)
	err := c.post(uri, in, out)
	return out, err
}

// CronDisable disables a cronjob.
func (c *client) CronUpdate(owner, name, cron string, in *CronPatch) (*Cron, error) {
	out := new(Cron)
	uri := fmt.Sprintf(pathCron, c.addr, owner, name, cron)
	err := c.patch(uri, in, out)
	return out, err
}

// CronDelete deletes a cronjob.
func (c *client) CronDelete(owner, name, cron string) error {
	uri := fmt.Sprintf(pathCron, c.addr, owner, name, cron)
	return c.delete(uri)
}

// Queue returns a list of enqueued builds.
func (c *client) Queue() ([]*Stage, error) {
	var out []*Stage
	uri := fmt.Sprintf(pathQueue, c.addr)
	err := c.get(uri, &out)
	return out, err
}

// QueueResume resumes queue operations.
func (c *client) QueueResume() error {
	uri := fmt.Sprintf(pathQueue, c.addr)
	err := c.post(uri, nil, nil)
	return err
}

// QueuePause pauses queue operations.
func (c *client) QueuePause() error {
	uri := fmt.Sprintf(pathQueue, c.addr)
	err := c.delete(uri)
	return err
}

// Node returns a node by name.
func (c *client) Node(name string) (*Node, error) {
	out := new(Node)
	uri := fmt.Sprintf(pathNode, c.addr, name)
	err := c.get(uri, out)
	return out, err
}

// NodeList returns a list of all nodes.
func (c *client) NodeList() ([]*Node, error) {
	var out []*Node
	uri := fmt.Sprintf(pathNodes, c.addr)
	err := c.get(uri, &out)
	return out, err
}

// NodeCreate creates a node.
func (c *client) NodeCreate(in *Node) (*Node, error) {
	out := new(Node)
	uri := fmt.Sprintf(pathNodes, c.addr)
	err := c.post(uri, in, out)
	return out, err
}

// NodeDelete deletes a node.
func (c *client) NodeDelete(name string) error {
	uri := fmt.Sprintf(pathNode, c.addr, name)
	return c.delete(uri)
}

// NodeUpdate updates a node.
func (c *client) NodeUpdate(name string, in *NodePatch) (*Node, error) {
	out := new(Node)
	uri := fmt.Sprintf(pathNode, c.addr, name)
	err := c.patch(uri, in, out)
	return out, err
}

//
// autoscaler
//

// Server returns the named servers details.
func (c *client) Server(name string) (*Server, error) {
	out := new(Server)
	uri := fmt.Sprintf(pathServer, c.addr, name)
	err := c.get(uri, &out)
	return out, err
}

// ServerList returns a list of all active build servers.
func (c *client) ServerList() ([]*Server, error) {
	var out []*Server
	uri := fmt.Sprintf(pathServers, c.addr)
	err := c.get(uri, &out)
	return out, err
}

// ServerCreate creates a new server.
func (c *client) ServerCreate() (*Server, error) {
	out := new(Server)
	uri := fmt.Sprintf(pathServers, c.addr)
	err := c.post(uri, nil, out)
	return out, err
}

// ServerDelete terminates a server.
func (c *client) ServerDelete(name string) error {
	uri := fmt.Sprintf(pathServer, c.addr, name)
	return c.delete(uri)
}

// AutoscalePause pauses the autoscaler.
func (c *client) AutoscalePause() error {
	uri := fmt.Sprintf(pathScalerPause, c.addr)
	return c.post(uri, nil, nil)
}

// AutoscaleResume resumes the autoscaler.
func (c *client) AutoscaleResume() error {
	uri := fmt.Sprintf(pathScalerResume, c.addr)
	return c.post(uri, nil, nil)
}

// AutoscaleVersion resumes the autoscaler.
func (c *client) AutoscaleVersion() (*Version, error) {
	out := new(Version)
	uri := fmt.Sprintf(pathVersion, c.addr)
	err := c.get(uri, out)
	return out, err
}

//
// http request helper functions
//

// helper function for making an http GET request.
func (c *client) get(rawurl string, out interface{}) error {
	return c.do(rawurl, "GET", nil, out)
}

// helper function for making an http POST request.
func (c *client) post(rawurl string, in, out interface{}) error {
	return c.do(rawurl, "POST", in, out)
}

// helper function for making an http PATCH request.
func (c *client) patch(rawurl string, in, out interface{}) error {
	return c.do(rawurl, "PATCH", in, out)
}

// helper function for making an http DELETE request.
func (c *client) delete(rawurl string) error {
	return c.do(rawurl, "DELETE", nil, nil)
}

// helper function to make an http request
func (c *client) do(rawurl, method string, in, out interface{}) error {
	body, err := c.open(rawurl, method, in, out)
	if err != nil {
		return err
	}
	defer body.Close()
	if out != nil {
		return json.NewDecoder(body).Decode(out)
	}
	return nil
}

// helper function to open an http request
func (c *client) open(rawurl, method string, in, out interface{}) (io.ReadCloser, error) {
	uri, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, uri.String(), nil)
	if err != nil {
		return nil, err
	}
	if in != nil {
		decoded, derr := json.Marshal(in)
		if derr != nil {
			return nil, derr
		}
		buf := bytes.NewBuffer(decoded)
		req.Body = ioutil.NopCloser(buf)
		req.ContentLength = int64(len(decoded))
		req.Header.Set("Content-Length", strconv.Itoa(len(decoded)))
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 299 {
		defer resp.Body.Close()
		out, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("client error %d: %s", resp.StatusCode, string(out))
	}
	return resp.Body, nil
}

// mapValues converts a map to url.Values
func mapValues(params map[string]string) url.Values {
	values := url.Values{}
	for key, val := range params {
		values.Add(key, val)
	}
	return values
}
