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

// Architectures
const (
	Freebsd_386 uint = iota
	Freebsd_amd64
	Freebsd_arm
	Linux_386
	Linux_amd64
	Linux_arm
	Linux_arm64
	Solaris_amd64
	Windows_386
	Windows_amd64
)

// Architecture Map
var Archs = map[string]uint{
	"freebsd_386":   Freebsd_386,
	"freebsd_amd64": Freebsd_amd64,
	"freebsd_arm":   Freebsd_arm,
	"linux_386":     Linux_386,
	"linux_amd64":   Linux_amd64,
	"linux_arm":     Linux_arm,
	"linux_arm64":   Linux_arm64,
	"solaris_amd64": Solaris_amd64,
	"windows_386":   Windows_386,
	"windows_amd64": Windows_amd64,
}
