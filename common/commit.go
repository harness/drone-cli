package common

type Commit struct {
	Status      string `json:"status"`
	Started     int64  `json:"started_at"`
	Finished    int64  `json:"finished_at"`
	Duration    int64  `json:"duration"`
	Sha         string `json:"sha"`
	Branch      string `json:"branch"`
	PullRequest string `json:"pull_request"`
	Author      string `json:"author"`
	Gravatar    string `json:"gravatar"`
	Timestamp   string `json:"timestamp"`
	Message     string `json:"message"`
}
