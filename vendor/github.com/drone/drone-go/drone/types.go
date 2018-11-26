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

type (
	// User represents a user account.
	User struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Email     string `json:"email"`
		Avatar    string `json:"avatar_url"`
		Active    bool   `json:"active"`
		Admin     bool   `json:"admin"`
		Machine   bool   `json:"machine"`
		Syncing   bool   `json:"syncing"`
		Synced    int64  `json:"synced"`
		Created   int64  `json:"created"`
		Updated   int64  `json:"updated"`
		LastLogin int64  `json:"last_login"`
	}

	// Repo represents a repository.
	Repo struct {
		ID         int64  `json:"id"`
		UID        string `json:"uid"`
		UserID     int64  `json:"user_id"`
		Namespace  string `json:"namespace"`
		Name       string `json:"name"`
		Slug       string `json:"slug"`
		SCM        string `json:"scm"`
		HTTPURL    string `json:"git_http_url"`
		SSHURL     string `json:"git_ssh_url"`
		Link       string `json:"link"`
		Branch     string `json:"default_branch"`
		Private    bool   `json:"private"`
		Visibility string `json:"visibility"`
		Active     bool   `json:"active"`
		Config     string `json:"config_path"`
		Trusted    bool   `json:"trusted"`
		Protected  bool   `json:"protected"`
		Timeout    int64  `json:"timeout"`
		Counter    int64  `json:"counter"`
		Synced     int64  `json:"synced"`
		Created    int64  `json:"created"`
		Updated    int64  `json:"updated"`
		Version    int64  `json:"version"`
	}

	// RepoPatch defines a repository patch request.
	RepoPatch struct {
		Config     *string `json:"config_path,omitempty"`
		Protected  *bool   `json:"protected,omitempty"`
		Trusted    *bool   `json:"trusted,omitempty"`
		Timeout    *int64  `json:"timeout,omitempty"`
		Visibility *string `json:"visibility,omitempty"`
		Counter    *int    `json:"counter,omitempty"`
	}

	// Build defines a build object.
	Build struct {
		ID           int64             `json:"id"`
		RepoID       int64             `json:"repo_id"`
		Trigger      string            `json:"trigger"`
		Number       int64             `json:"number"`
		Parent       int64             `json:"parent,omitempty"`
		Status       string            `json:"status"`
		Error        string            `json:"error,omitempty"`
		Event        string            `json:"event"`
		Action       string            `json:"action"`
		Link         string            `json:"link"`
		Timestamp    int64             `json:"timestamp"`
		Title        string            `json:"title,omitempty"`
		Message      string            `json:"message"`
		Before       string            `json:"before"`
		After        string            `json:"after"`
		Ref          string            `json:"ref"`
		Fork         string            `json:"source_repo"`
		Source       string            `json:"source"`
		Target       string            `json:"target"`
		Author       string            `json:"author_login"`
		AuthorName   string            `json:"author_name"`
		AuthorEmail  string            `json:"author_email"`
		AuthorAvatar string            `json:"author_avatar"`
		Sender       string            `json:"sender"`
		Params       map[string]string `json:"params,omitempty"`
		Deploy       string            `json:"deploy_to,omitempty"`
		Started      int64             `json:"started"`
		Finished     int64             `json:"finished"`
		Created      int64             `json:"created"`
		Updated      int64             `json:"updated"`
		Version      int64             `json:"version"`
		Stages       []*Stage          `json:"stages,omitempty"`
	}

	// Stage represents a stage of build execution.
	Stage struct {
		ID        int64             `json:"id"`
		BuildID   int64             `json:"build_id"`
		Number    int               `json:"number"`
		Name      string            `json:"name"`
		Kind      string            `json:"kind,omitempty"`
		Type      string            `json:"type,omitempty"`
		Status    string            `json:"status"`
		Error     string            `json:"error,omitempty"`
		ErrIgnore bool              `json:"errignore"`
		ExitCode  int               `json:"exit_code"`
		Machine   string            `json:"machine,omitempty"`
		OS        string            `json:"os"`
		Arch      string            `json:"arch"`
		Variant   string            `json:"variant,omitempty"`
		Kernel    string            `json:"kernel,omitempty"`
		Started   int64             `json:"started"`
		Stopped   int64             `json:"stopped"`
		Created   int64             `json:"created"`
		Updated   int64             `json:"updated"`
		Version   int64             `json:"version"`
		OnSuccess bool              `json:"on_success"`
		OnFailure bool              `json:"on_failure"`
		DependsOn []string          `json:"depends_on,omitempty"`
		Labels    map[string]string `json:"labels,omitempty"`
		Steps     []*Step           `json:"steps,omitempty"`
	}

	// Step represents an individual step in the stage.
	Step struct {
		ID        int64  `json:"id"`
		StageID   int64  `json:"step_id"`
		Number    int    `json:"number"`
		Name      string `json:"name"`
		Status    string `json:"status"`
		Error     string `json:"error,omitempty"`
		ErrIgnore bool   `json:"errignore,omitempty"`
		ExitCode  int    `json:"exit_code"`
		Started   int64  `json:"started,omitempty"`
		Stopped   int64  `json:"stopped,omitempty"`
		Version   int64  `json:"version"`
	}

	// Registry represents a docker registry with credentials.
	// DEPRECATED
	Registry struct {
		Address  string `json:"address"`
		Username string `json:"username"`
		Password string `json:"password,omitempty"`
		Email    string `json:"email"`
		Token    string `json:"token"`
		Policy   string `json:"policy,omitempty"`
	}

	// Secret represents a secret variable, such as a password or token.
	Secret struct {
		Name            string `json:"name,omitempty"`
		Data            string `json:"data,omitempty"`
		PullRequest     bool   `json:"pull_request,omitempty"`
		PullRequestPush bool   `json:"pull_request_push,omitempty"`

		// Deprecated.
		Pull bool `json:"pull,omitempty"`
		Fork bool `json:"fork,omitempty"`
	}

	// Server represents a server node.
	Server struct {
		ID       string `json:"id"`
		Provider string `json:"provider"`
		State    string `json:"state"`
		Name     string `json:"name"`
		Image    string `json:"image"`
		Region   string `json:"region"`
		Size     string `json:"size"`
		Address  string `json:"address"`
		Capacity int    `json:"capacity"`
		Secret   string `json:"secret"`
		Error    string `json:"error"`
		CAKey    []byte `json:"ca_key"`
		CACert   []byte `json:"ca_cert"`
		TLSKey   []byte `json:"tls_key"`
		TLSCert  []byte `json:"tls_cert"`
		Created  int64  `json:"created"`
		Updated  int64  `json:"updated"`
		Started  int64  `json:"started"`
		Stopped  int64  `json:"stopped"`
	}

	// Cron represents a cron job.
	Cron struct {
		ID       int64  `json:"id"`
		RepoID   int64  `json:"repo_id"`
		Name     string `json:"name"`
		Expr     string `json:"expr"`
		Next     int64  `json:"next"`
		Prev     int64  `json:"prev"`
		Event    string `json:"event"`
		Branch   string `json:"branch"`
		Target   string `json:"target"`
		Disabled bool   `json:"disabled"`
		Created  int64  `json:"created"`
		Updated  int64  `json:"updated"`
	}

	// CronPatch defines a cron patch request.
	CronPatch struct {
		Event    *string `json:"event"`
		Branch   *string `json:"branch"`
		Target   *string `json:"target"`
		Disabled *bool   `json:"disabled"`
	}

	// Line represents a line of container logs.
	Line struct {
		Number    int    `json:"pos"`
		Message   string `json:"out"`
		Timestamp int64  `json:"time"`
	}

	// Config represents a config file.
	Config struct {
		Data string `json:"data"`
		Kind string `json:"kind"`
	}

	// Version provides system version details.
	Version struct {
		Source  string `json:"source,omitempty"`
		Version string `json:"version,omitempty"`
		Commit  string `json:"commit,omitempty"`
	}

	// Node provides node details.
	Node struct {
		ID        int64             `json:"id"`
		UID       string            `json:"uid"`
		Provider  string            `json:"provider"`
		State     string            `json:"state"`
		Name      string            `json:"name"`
		Image     string            `json:"image"`
		Region    string            `json:"region"`
		Size      string            `json:"size"`
		OS        string            `json:"os"`
		Arch      string            `json:"arch"`
		Kernel    string            `json:"kernel"`
		Variant   string            `json:"variant"`
		Address   string            `json:"address"`
		Capacity  int               `json:"capacity"`
		Filters   []string          `json:"filters"`
		Labels    map[string]string `json:"labels"`
		Error     string            `json:"error"`
		CAKey     []byte            `json:"ca_key"`
		CACert    []byte            `json:"ca_cert"`
		TLSKey    []byte            `json:"tls_key"`
		TLSCert   []byte            `json:"tls_cert"`
		TLSName   string            `json:"tls_name"`
		Paused    bool              `json:"paused"`
		Protected bool              `json:"protected"`
		Created   int64             `json:"created"`
		Updated   int64             `json:"updated"`
	}

	// NodePatch defines a node patch request.
	NodePatch struct {
		UID       *string            `json:"uid"`
		Provider  *string            `json:"provider"`
		State     *string            `json:"state"`
		Image     *string            `json:"image"`
		Region    *string            `json:"region"`
		Size      *string            `json:"size"`
		Address   *string            `json:"address"`
		Capacity  *int               `json:"capacity"`
		Filters   *[]string          `json:"filters"`
		Labels    *map[string]string `json:"labels"`
		Error     *string            `json:"error"`
		CAKey     *[]byte            `json:"ca_key"`
		CACert    *[]byte            `json:"ca_cert"`
		TLSKey    *[]byte            `json:"tls_key"`
		TLSCert   *[]byte            `json:"tls_cert"`
		Paused    *bool              `json:"paused"`
		Protected *bool              `json:"protected"`
	}
)

// Error represents a json-encoded API error.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}
