package drone

// User represents a user account.
type User struct {
	ID     int64  `json:"id"`
	Login  string `json:"login"`
	Email  string `json:"email"`
	Avatar string `json:"avatar_url"`
	Active bool   `json:"active"`
	Admin  bool   `json:"admin"`
}

// Repo represents a version control repository.
type Repo struct {
	ID          int64  `json:"id"`
	Owner       string `json:"owner"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Avatar      string `json:"avatar_url"`
	Link        string `json:"link_url"`
	Kind        string `json:"kind"`
	Clone       string `json:"clone_url"`
	Branch      string `json:"default_branch"`
	Timeout     int64  `json:"timeout"`
	IsPrivate   bool   `json:"private"`
	IsTrusted   bool   `json:"trusted"`
	AllowPull   bool   `json:"allow_pr"`
	AllowPush   bool   `json:"allow_push"`
	AllowDeploy bool   `json:"allow_deploys"`
	AllowTag    bool   `json:"allow_tags"`
}

// Build represents the process of compiling and testing a changeset, typically
// triggered by the remote system (ie GitHub).
type Build struct {
	ID        int64  `json:"id"`
	Number    int    `json:"number"`
	Event     string `json:"event"`
	Status    string `json:"status"`
	Enqueued  int64  `json:"enqueued_at"`
	Created   int64  `json:"created_at"`
	Started   int64  `json:"started_at"`
	Finished  int64  `json:"finished_at"`
	Deploy    string `json:"deploy_to"`
	Commit    string `json:"commit"`
	Branch    string `json:"branch"`
	Ref       string `json:"ref"`
	Refspec   string `json:"refspec"`
	Remote    string `json:"remote"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
	Author    string `json:"author"`
	Avatar    string `json:"author_avatar"`
	Email     string `json:"author_email"`
	Link      string `json:"link_url"`
	Signed    bool   `json:"signed"`
	Verified  bool   `json:"verified"`
}

// Job represents a single job that is being executed as part of a Build.
type Job struct {
	ID       int64  `json:"id"`
	Number   int    `json:"number"`
	Status   string `json:"status"`
	Error    string `json:"error"`
	ExitCode int    `json:"exit_code"`
	Enqueued int64  `json:"enqueued_at"`
	Started  int64  `json:"started_at"`
	Finished int64  `json:"finished_at"`

	Environment map[string]string `json:"environment"`
}

// Secret represents a repository secret.
type Secret struct {
	ID     int64    `json:"id"`
	Name   string   `json:"name"`
	Value  string   `json:"value"`
	Images []string `json:"image"`
	Events []string `json:"event"`
}

// Activity represents a build activity. It combines the build details with
// summary Repository information.
type Activity struct {
	Owner     string `json:"owner"`
	Name      string `json:"name"`
	FullName  string `json:"full_name"`
	Number    int    `json:"number"`
	Event     string `json:"event"`
	Status    string `json:"status"`
	Enqueued  int64  `json:"enqueued_at"`
	Created   int64  `json:"created_at"`
	Started   int64  `json:"started_at"`
	Finished  int64  `json:"finished_at"`
	Commit    string `json:"commit"`
	Branch    string `json:"branch"`
	Ref       string `json:"ref"`
	Refspec   string `json:"refspec"`
	Remote    string `json:"remote"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
	Author    string `json:"author"`
	Avatar    string `json:"author_avatar"`
	Email     string `json:"author_email"`
	Link      string `json:"link_url"`
	Signed    bool   `json:"signed"`
	Verified  bool   `json:"verified"`
}

// Netrc defines a default .netrc file that should be injected into the build
// environment. It will be used to authorize access to https resources, such as
// git+https clones.
type Netrc struct {
	Machine  string `json:"machine"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// System represents the drone system.
type System struct {
	Version string `json:"version"`
	Link    string `json:"link_url"`
}

// Agent represents a registered build node.
type Agent struct {
	ID       int64  `json:"id"`
	Address  string `json:"address"`
	Platform string `json:"platform"`
	Capacity int    `json:"capacity"`
	Created  int64  `json:"created_at"`
	Updated  int64  `json:"updated_at"`
}

// Payload defines the full payload send to plugins.
type Payload struct {
	Yaml      string    `json:"config"`
	User      *User     `json:"user"`
	Repo      *Repo     `json:"repo"`
	Build     *Build    `json:"build"`
	BuildLast *Build    `json:"build_last"`
	Job       *Job      `json:"job"`
	Netrc     *Netrc    `json:"netrc"`
	System    *System   `json:"system"`
	Secrets   []*Secret `json:"secrets"`
}
