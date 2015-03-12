package common

// Build represents a build
type Build struct {
	Status   string `json:"status"`
	ExitCode int    `json:"exit_code"`
	Started  int64  `json:"started_at"`
	Finished int64  `json:"finished_at"`
	Duration int64  `json:"duration"`
	Label    string `json:"label"`
}
