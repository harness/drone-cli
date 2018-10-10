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
	"net/http"
)

// TODO(bradrydzewski) add repo + latest build endpoint
// TODO(bradrydzewski) add queue endpoint
// TDOO(bradrydzewski) add stats endpoint
// TODO(bradrydzewski) add version endpoint

// Client is used to communicate with a Drone server.
type Client interface {
	// SetClient sets the http.Client.
	SetClient(*http.Client)

	// SetAddress sets the server address.
	SetAddress(string)

	// Self returns the currently authenticated user.
	Self() (*User, error)

	// User returns a user by login.
	User(login string) (*User, error)

	// UserList returns a list of all registered users.
	UserList() ([]*User, error)

	// UserCreate creates a new user account.
	UserCreate(user *User) (*User, error)

	// UserUpdate updates a user account.
	UserUpdate(user *User) (*User, error)

	// UserDelete deletes a user account.
	UserDelete(login string) error

	// Repo returns a repository by name.
	Repo(namespace, name string) (*Repo, error)

	// RepoList returns a list of all repositories to which
	// the user has explicit access in the host system.
	RepoList() ([]*Repo, error)

	// RepoListSync returns a list of all repositories to which
	// the user has explicit access in the host system.
	RepoListSync() ([]*Repo, error)

	// RepoEnable activates a repository.
	RepoEnable(namespace, name string) (*Repo, error)

	// RepoUpdate updates a repository.
	RepoUpdate(namespace, name string, repo *RepoPatch) (*Repo, error)

	// RepoChown updates a repository owner.
	RepoChown(namespace, name string) (*Repo, error)

	// RepoRepair repairs the repository hooks.
	RepoRepair(namespace, name string) error

	// RepoDisable disables a repository.
	RepoDisable(namespace, name string) error

	// Build returns a repository build by number.
	Build(namespace, name string, build int) (*Build, error)

	// BuildLast returns the latest build by branch. An
	// empty branch will result in the default branch.
	BuildLast(namespace, name, branch string) (*Build, error)

	// BuildList returns a list of recent builds for the
	// the specified repository.
	BuildList(namespace, name string) ([]*Build, error)

	// BuildQueue returns a list of enqueued builds.
	BuildQueue() ([]*Build, error)

	// BuildRestart re-starts a build.
	BuildRestart(namespace, name string, build int, params map[string]string) (*Build, error)

	// BuildCancel stops the specified running job for
	// given build.
	BuildCancel(namespace, name string, build int) error

	// Approve approves a blocked build stage.
	Approve(namespace, name string, build, stage int) error

	// Decline declines a blocked build stage.
	Decline(namespace, name string, build, stage int) error

	// Promote promotes a build to the target environment.
	Promote(namespace, name string, build int, target string, params map[string]string) (*Build, error)

	// Rollback reverts the target environment to an previous build.
	Rollback(namespace, name string, build int, target string, params map[string]string) (*Build, error)

	// Logs gets the logs for the specified step.
	Logs(owner, name string, build, stage, step int) ([]*Line, error)

	// LogsPurge purges the build logs for the specified step.
	LogsPurge(owner, name string, build, stage, step int) error

	// Cron returns a cronjob by name.
	Cron(owner, name, cron string) (*Cron, error)

	// CronList returns a list of all repository cronjobs.
	CronList(owner string, name string) ([]*Cron, error)

	// CronEnable enables a cronjob.
	CronEnable(owner, name, cron string) error

	// CronDisable disables a cronjob.
	CronDisable(owner, name, cron string) error

	// Sign signs the yaml file.
	Sign(owner, name, file string) (string, error)

	// Verify verifies the yaml signature.
	Verify(owner, name, file string) error

	// Encrypt returns an encrypted secret
	Encrypt(owner, name string, secret *Secret) (string, error)

	//
	// Move to autoscaler-go
	//

	// Server returns the named servers details.
	Server(name string) (*Server, error)

	// ServerList returns a list of all active build servers.
	ServerList() ([]*Server, error)

	// ServerCreate creates a new server.
	ServerCreate() (*Server, error)

	// ServerDelete terminates a server.
	ServerDelete(name string) error

	// AutoscalePause pauses the autoscaler.
	AutoscalePause() error

	// AutoscaleResume resumes the autoscaler.
	AutoscaleResume() error

	// AutoscaleVersion returns the autoscaler version.
	AutoscaleVersion() (*Version, error)
}
