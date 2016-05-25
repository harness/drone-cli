package drone

// Events
const (
	EventPush   = "push"
	EventPull   = "pull_request"
	EventTag    = "tag"
	EventDeploy = "deployment"
)

// Statuses
const (
	StatusSkipped = "skipped"
	StatusPending = "pending"
	StatusRunning = "running"
	StatusSuccess = "success"
	StatusFailure = "failure"
	StatusKilled  = "killed"
	StatusError   = "error"
)
