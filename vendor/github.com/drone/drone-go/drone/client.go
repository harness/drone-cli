package drone

//go:generate mockery -all
//go:generate mv mocks/Client.go mocks/client.go

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

const (
	pathSelf    = "%s/api/user"
	pathFeed    = "%s/api/user/feed"
	pathRepos   = "%s/api/user/repos"
	pathRepo    = "%s/api/repos/%s/%s"
	pathEncrypt = "%s/api/repos/%s/%s/encrypt"
	pathBuilds  = "%s/api/repos/%s/%s/builds"
	pathBuild   = "%s/api/repos/%s/%s/builds/%v"
	pathJob     = "%s/api/repos/%s/%s/builds/%d/%d"
	pathLog     = "%s/api/repos/%s/%s/logs/%d/%d"
	pathKey     = "%s/api/repos/%s/%s/key"
	pathSign    = "%s/api/repos/%s/%s/sign"
	pathSecrets = "%s/api/repos/%s/%s/secrets"
	pathSecret  = "%s/api/repos/%s/%s/secrets/%s"
	pathNodes   = "%s/api/nodes"
	pathNode    = "%s/api/nodes/%d"
	pathUsers   = "%s/api/users"
	pathUser    = "%s/api/users/%s"
)

type client struct {
	client *http.Client
	base   string // base url
}

// NewClient returns a client at the specified url.
func NewClient(uri string) Client {
	return &client{http.DefaultClient, uri}
}

// NewClientToken returns a client at the specified url that
// authenticates all outbound requests with the given token.
func NewClientToken(uri, token string) Client {
	config := new(oauth2.Config)
	auther := config.Client(oauth2.NoContext, &oauth2.Token{AccessToken: token})
	return &client{auther, uri}
}

// NewClientTokenTLS returns a client at the specified url that
// authenticates all outbound requests with the given token and
// tls.Config if provided.
func NewClientTokenTLS(uri, token string, c *tls.Config) Client {
	config := new(oauth2.Config)
	auther := config.Client(oauth2.NoContext, &oauth2.Token{AccessToken: token})
	if c != nil {
		auther.Transport.(*oauth2.Transport).Base = &http.Transport{TLSClientConfig: c}
	}
	return &client{auther, uri}
}

// SetClient sets the default http client. This should be
// used in conjunction with golang.org/x/oauth2 to
// authenticate requests to the Drone server.
func (c *client) SetClient(client *http.Client) {
	c.client = client
}

// Self returns the currently authenticated user.
func (c *client) Self() (*User, error) {
	out := new(User)
	uri := fmt.Sprintf(pathSelf, c.base)
	err := c.get(uri, out)
	return out, err
}

// User returns a user by login.
func (c *client) User(login string) (*User, error) {
	out := new(User)
	uri := fmt.Sprintf(pathUser, c.base, login)
	err := c.get(uri, out)
	return out, err
}

// UserList returns a list of all registered users.
func (c *client) UserList() ([]*User, error) {
	var out []*User
	uri := fmt.Sprintf(pathUsers, c.base)
	err := c.get(uri, &out)
	return out, err
}

// UserPost creates a new user account.
func (c *client) UserPost(in *User) (*User, error) {
	out := new(User)
	uri := fmt.Sprintf(pathUsers, c.base)
	err := c.post(uri, in, out)
	return out, err
}

// UserPatch updates a user account.
func (c *client) UserPatch(in *User) (*User, error) {
	out := new(User)
	uri := fmt.Sprintf(pathUser, c.base, in.Login)
	err := c.patch(uri, in, out)
	return out, err
}

// UserDel deletes a user account.
func (c *client) UserDel(login string) error {
	uri := fmt.Sprintf(pathUser, c.base, login)
	err := c.delete(uri)
	return err
}

// UserFeed returns the user's activity feed.
func (c *client) UserFeed() ([]*Activity, error) {
	var out []*Activity
	uri := fmt.Sprintf(pathFeed, c.base)
	err := c.get(uri, &out)
	return out, err
}

// Repo returns a repository by name.
func (c *client) Repo(owner string, name string) (*Repo, error) {
	out := new(Repo)
	uri := fmt.Sprintf(pathRepo, c.base, owner, name)
	err := c.get(uri, out)
	return out, err
}

// RepoList returns a list of all repositories to which
// the user has explicit access in the host system.
func (c *client) RepoList() ([]*Repo, error) {
	var out []*Repo
	uri := fmt.Sprintf(pathRepos, c.base)
	err := c.get(uri, &out)
	return out, err
}

// RepoPost activates a repository.
func (c *client) RepoPost(owner string, name string) (*Repo, error) {
	out := new(Repo)
	uri := fmt.Sprintf(pathRepo, c.base, owner, name)
	err := c.post(uri, nil, out)
	return out, err
}

// RepoPatch updates a repository.
func (c *client) RepoPatch(in *Repo) (*Repo, error) {
	out := new(Repo)
	uri := fmt.Sprintf(pathRepo, c.base, in.Owner, in.Name)
	err := c.patch(uri, in, out)
	return out, err
}

// RepoDel deletes a repository.
func (c *client) RepoDel(owner, name string) error {
	uri := fmt.Sprintf(pathRepo, c.base, owner, name)
	err := c.delete(uri)
	return err
}

// RepoKey returns a repository public key.
func (c *client) RepoKey(owner, name string) (*Key, error) {
	out := new(Key)
	uri := fmt.Sprintf(pathKey, c.base, owner, name)
	rc, err := c.stream(uri, "GET", nil, nil)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	raw, _ := ioutil.ReadAll(rc)
	out.Public = string(raw)
	return out, err
}

// Build returns a repository build by number.
func (c *client) Build(owner, name string, num int) (*Build, error) {
	out := new(Build)
	uri := fmt.Sprintf(pathBuild, c.base, owner, name, num)
	err := c.get(uri, out)
	return out, err
}

// Build returns the latest repository build by branch.
func (c *client) BuildLast(owner, name, branch string) (*Build, error) {
	out := new(Build)
	uri := fmt.Sprintf(pathBuild, c.base, owner, name, "latest")
	if len(branch) != 0 {
		uri += "?branch=" + branch
	}
	err := c.get(uri, out)
	return out, err
}

// BuildList returns a list of recent builds for the
// the specified repository.
func (c *client) BuildList(owner, name string) ([]*Build, error) {
	var out []*Build
	uri := fmt.Sprintf(pathBuilds, c.base, owner, name)
	err := c.get(uri, &out)
	return out, err
}

// BuildStart re-starts a stopped build.
func (c *client) BuildStart(owner, name string, num int) (*Build, error) {
	out := new(Build)
	uri := fmt.Sprintf(pathBuild, c.base, owner, name, num)
	err := c.post(uri, nil, out)
	return out, err
}

// BuildStop cancels the running job.
func (c *client) BuildStop(owner, name string, num, job int) error {
	uri := fmt.Sprintf(pathJob, c.base, owner, name, num, job)
	err := c.delete(uri)
	return err
}

// BuildFork re-starts a stopped build with a new build number,
// preserving the prior history.
func (c *client) BuildFork(owner, name string, num int) (*Build, error) {
	out := new(Build)
	uri := fmt.Sprintf(pathBuild+"?fork=true", c.base, owner, name, num)
	err := c.post(uri, nil, out)
	return out, err
}

// BuildLogs returns the build logs for the specified job.
func (c *client) BuildLogs(owner, name string, num, job int) (io.ReadCloser, error) {
	uri := fmt.Sprintf(pathLog, c.base, owner, name, num, job)
	return c.stream(uri, "GET", nil, nil)
}

// Deploy triggers a deployment for an existing build using the
// specified target environment.
func (c *client) Deploy(owner, name string, num int, env string) (*Build, error) {
	out := new(Build)
	val := url.Values{}
	val.Set("fork", "true")
	val.Set("event", "deployment")
	val.Set("deploy_to", env)
	uri := fmt.Sprintf(pathBuild+"?"+val.Encode(), c.base, owner, name, num)
	err := c.post(uri, nil, out)
	return out, err
}

// SecretPost create or updates a repository secret.
func (c *client) SecretPost(owner, name string, secret *Secret) error {
	uri := fmt.Sprintf(pathSecrets, c.base, owner, name)
	return c.post(uri, secret, nil)
}

// SecretDel deletes a named repository secret.
func (c *client) SecretDel(owner, name, secret string) error {
	uri := fmt.Sprintf(pathSecret, c.base, owner, name, secret)
	return c.delete(uri)
}

// Sign returns a cryptographic signature for the input string.
func (c *client) Sign(owner, name string, in []byte) ([]byte, error) {
	buf := bytes.Buffer{}
	buf.Write(in)
	uri := fmt.Sprintf(pathSign, c.base, owner, name)
	rc, err := c.stream(uri, "POST", buf, nil)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return ioutil.ReadAll(rc)
}

// Node returns a node by id.
func (c *client) Node(id int64) (*Node, error) {
	out := new(Node)
	uri := fmt.Sprintf(pathNode, c.base, id)
	err := c.get(uri, out)
	return out, err
}

// NodeList returns a list of all registered worker nodes.
func (c *client) NodeList() ([]*Node, error) {
	var out []*Node
	uri := fmt.Sprintf(pathNodes, c.base)
	err := c.get(uri, &out)
	return out, err
}

// NodePost registers a new worker node.
func (c *client) NodePost(in *Node) (*Node, error) {
	out := new(Node)
	uri := fmt.Sprintf(pathNodes, c.base)
	err := c.post(uri, in, out)
	return out, err
}

// NodeDel deletes a worker node.
func (c *client) NodeDel(id int64) error {
	uri := fmt.Sprintf(pathNode, c.base, id)
	err := c.delete(uri)
	return err
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

// helper function for making an http PUT request.
func (c *client) put(rawurl string, in, out interface{}) error {
	return c.do(rawurl, "PUT", in, out)
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
	// executes the http request and returns the body as
	// and io.ReadCloser
	body, err := c.stream(rawurl, method, in, out)
	if err != nil {
		return err
	}
	defer body.Close()

	// if a json response is expected, parse and return
	// the json response.
	if out != nil {
		return json.NewDecoder(body).Decode(out)
	}
	return nil
}

// helper function to stream an http request
func (c *client) stream(rawurl, method string, in, out interface{}) (io.ReadCloser, error) {
	uri, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	// if we are posting or putting data, we need to
	// write it to the body of the request.
	var buf io.ReadWriter
	if in == nil {
		// nothing
	} else if rw, ok := in.(io.ReadWriter); ok {
		buf = rw
	} else {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(in)
		if err != nil {
			return nil, err
		}
	}

	// creates a new http request to bitbucket.
	req, err := http.NewRequest(method, uri.String(), buf)
	if err != nil {
		return nil, err
	}
	if in == nil {
		// nothing
	} else if _, ok := in.(io.ReadWriter); ok {
		req.Header.Set("Content-Type", "plain/text")
	} else {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > http.StatusPartialContent {
		defer resp.Body.Close()
		out, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf(string(out))
	}
	return resp.Body, nil
}
