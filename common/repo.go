package common

type Repo struct {
	ID     int64  `json:"id"`
	Remote string `json:"remote"`
	Host   string `json:"host"`
	Owner  string `json:"owner"`
	Name   string `json:"name"`
	URL    string `json:"url"`
}
