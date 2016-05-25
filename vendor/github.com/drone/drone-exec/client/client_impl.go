package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/drone/drone-go/drone"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
	"golang.org/x/oauth2"
)

const (
	pathPull     = "%s/api/queue/pull/%s/%s"
	pathWait     = "%s/api/queue/wait/%d"
	pathStream   = "%s/api/queue/stream/%d"
	pathPush     = "%s/api/queue/status/%d"
	pathPing     = "%s/api/queue/ping"
	pathLogs     = "%s/api/queue/logs/%d"
	pathLogsAuth = "%s/api/queue/logs/%d?access_token=%s"

	pathSelf       = "%s/api/user"
	pathFeed       = "%s/api/user/feed"
	pathRepos      = "%s/api/user/repos"
	pathRepo       = "%s/api/repos/%s/%s"
	pathEncrypt    = "%s/api/repos/%s/%s/encrypt"
	pathBuilds     = "%s/api/repos/%s/%s/builds"
	pathBuild      = "%s/api/repos/%s/%s/builds/%v"
	pathJob        = "%s/api/repos/%s/%s/builds/%d/%d"
	pathLog        = "%s/api/repos/%s/%s/logs/%d/%d"
	pathKey        = "%s/api/repos/%s/%s/key"
	pathSign       = "%s/api/repos/%s/%s/sign"
	pathSecrets    = "%s/api/repos/%s/%s/secrets"
	pathSecret     = "%s/api/repos/%s/%s/secrets/%s"
	pathNodes      = "%s/api/nodes"
	pathNode       = "%s/api/nodes/%d"
	pathUsers      = "%s/api/users"
	pathUser       = "%s/api/users/%s"
	pathBuildQueue = "%s/api/builds"
	pathAgent      = "%s/api/agents"
)

type client struct {
	client *http.Client
	token  string // auth token
	base   string // base url
}

// NewClient returns a client at the specified url.
func NewClient(uri string) Client {
	return &client{client: http.DefaultClient, base: uri}
}

// NewClientToken returns a client at the specified url that authenticates all
// outbound requests with the given token.
func NewClientToken(uri, token string) Client {
	config := new(oauth2.Config)
	auther := config.Client(oauth2.NoContext, &oauth2.Token{AccessToken: token})
	return &client{client: auther, base: uri, token: token}
}

// NewClientTokenTLS returns a client at the specified url that authenticates
// all outbound requests with the given token and tls.Config if provided.
func NewClientTokenTLS(uri, token string, c *tls.Config) Client {
	config := new(oauth2.Config)
	auther := config.Client(oauth2.NoContext, &oauth2.Token{AccessToken: token})
	if c != nil {
		if trans, ok := auther.Transport.(*oauth2.Transport); ok {
			trans.Base = &http.Transport{TLSClientConfig: c}
		}
	}
	return &client{client: auther, base: uri, token: token}
}

//
// below items for Queue (internal use only)
//

// Pull pulls work from the server queue.
func (c *client) Pull(os, arch string) (*drone.Payload, error) {
	out := new(drone.Payload)
	uri := fmt.Sprintf(pathPull, c.base, os, arch)
	err := c.post(uri, nil, out)
	return out, err
}

// Push pushes an update to the server.
func (c *client) Push(p *drone.Payload) error {
	uri := fmt.Sprintf(pathPush, c.base, p.Job.ID)
	err := c.post(uri, p, nil)
	return err
}

// Ping pings the server.
func (c *client) Ping() error {
	uri := fmt.Sprintf(pathPing, c.base)
	err := c.post(uri, nil, nil)
	return err
}

// Stream streams the build logs to the server.
func (c *client) Stream(id int64, rc io.ReadCloser) error {
	uri := fmt.Sprintf(pathStream, c.base, id)
	err := c.post(uri, rc, nil)

	return err
}

// LogPost sends the full build logs to the server.
func (c *client) LogPost(id int64, rc io.ReadCloser) error {
	uri := fmt.Sprintf(pathLogs, c.base, id)
	return c.post(uri, rc, nil)
}

// StreamWriter implements a special writer for streaming log entries to the
// central Drone server. The standard implementation is the gorilla.Connection.
type StreamWriter interface {
	Close() error
	WriteJSON(interface{}) error
}

// LogStream streams the build logs to the server.
func (c *client) LogStream(id int64) (StreamWriter, error) {
	rawurl := fmt.Sprintf(pathLogsAuth, c.base, id, c.token)
	uri, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	if uri.Scheme == "https" {
		uri.Scheme = "wss"
	} else {
		uri.Scheme = "ws"
	}

	// TODO need TLS client configuration

	conn, _, err := websocket.DefaultDialer.Dial(uri.String(), nil)
	return conn, err
}

// Wait watches and waits for the build to cancel or finish.
func (c *client) Wait(id int64) *Wait {
	ctx, cancel := context.WithCancel(context.Background())
	return &Wait{id, c, ctx, cancel}
}

type Wait struct {
	id     int64
	client *client

	ctx    context.Context
	cancel context.CancelFunc
}

func (w *Wait) Done() (*drone.Job, error) {
	uri := fmt.Sprintf(pathWait, w.client.base, w.id)
	req, err := w.client.createRequest(uri, "POST", nil)
	if err != nil {
		return nil, err
	}

	res, err := ctxhttp.Do(w.ctx, w.client.client, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	job := &drone.Job{}
	err = json.NewDecoder(res.Body).Decode(&job)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (w *Wait) Cancel() {
	w.cancel()
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
	body, err := c.open(rawurl, method, in, out)
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

// helper function to open an http request
func (c *client) open(rawurl, method string, in, out interface{}) (io.ReadCloser, error) {
	uri, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	// creates a new http request to bitbucket.
	req, err := http.NewRequest(method, uri.String(), nil)
	if err != nil {
		return nil, err
	}

	// if we are posting or putting data, we need to
	// write it to the body of the request.
	if in != nil {
		rc, ok := in.(io.ReadCloser)
		if ok {
			req.Body = rc
			req.Header.Set("Content-Type", "plain/text")
		} else {
			inJson, err := json.Marshal(in)
			if err != nil {
				return nil, err
			}

			buf := bytes.NewBuffer(inJson)
			req.Body = ioutil.NopCloser(buf)

			req.ContentLength = int64(len(inJson))
			req.Header.Set("Content-Length", strconv.Itoa(len(inJson)))
			req.Header.Set("Content-Type", "application/json")
		}
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > http.StatusPartialContent {
		defer resp.Body.Close()
		out, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("client error %d: %s", resp.StatusCode, string(out))
	}
	return resp.Body, nil
}

// createRequest is a helper function that builds an http.Request.
func (c *client) createRequest(rawurl, method string, in interface{}) (*http.Request, error) {
	uri, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	// if we are posting or putting data, we need to
	// write it to the body of the request.
	var buf io.ReadWriter
	if in != nil {
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
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}
