package drone

import "io"

// Client describes a drone client.
type Client interface {
	// Self returns the currently authenticated user.
	Self() (*User, error)

	// User returns a user by login.
	User(string) (*User, error)

	// UserList returns a list of all registered users.
	UserList() ([]*User, error)

	// UserPost creates a new user account.
	UserPost(*User) (*User, error)

	// UserPatch updates a user account.
	UserPatch(*User) (*User, error)

	// UserDel deletes a user account.
	UserDel(string) error

	// UserFeed returns the user's activity feed.
	UserFeed() ([]*Activity, error)

	// Repo returns a repository by name.
	Repo(string, string) (*Repo, error)

	// RepoList returns a list of all repositories to which the user has explicit
	// access in the host system.
	RepoList() ([]*Repo, error)

	// RepoPost activates a repository.
	RepoPost(string, string) (*Repo, error)

	// RepoPatch updates a repository.
	RepoPatch(*Repo) (*Repo, error)

	// RepoDel deletes a repository.
	RepoDel(string, string) error

	// Build returns a repository build by number.
	Build(string, string, int) (*Build, error)

	// BuildLast returns the latest repository build by branch. An empty branch
	// will result in the default branch.
	BuildLast(string, string, string) (*Build, error)

	// BuildList returns a list of recent builds for the
	// the specified repository.
	BuildList(string, string) ([]*Build, error)

	// BuildStart re-starts a stopped build.
	BuildStart(string, string, int) (*Build, error)

	// BuildStop stops the specified running job for given build.
	BuildStop(string, string, int, int) error

	// BuildFork re-starts a stopped build with a new build number, preserving
	// the prior history.
	BuildFork(string, string, int) (*Build, error)

	// BuildLogs returns the build logs for the specified job.
	BuildLogs(string, string, int, int) (io.ReadCloser, error)

	// BuildQueue returns a list of builds in queue.
	BuildQueue() ([]*Activity, error)

	// Deploy triggers a deployment for an existing build using the specified
	// target environment.
	Deploy(string, string, int, string) (*Build, error)

	// Sign returns a cryptographic signature for the input string.
	Sign(string, string, []byte) ([]byte, error)

	// SecretPost create or updates a repository secret.
	SecretPost(string, string, *Secret) error

	// SecretDel deletes a named repository secret.
	SecretDel(string, string, string) error

	// Agent returns an agent by IP address.
	Agent(string) (*Agent, error)

	// AgentDel deletes an agent by IP address.
	AgentDel(string) error

	// AgentList returns a list of build agents.
	AgentList() ([]*Agent, error)
}
