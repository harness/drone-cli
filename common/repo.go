package common

type Repo struct {
	Remote string `json:"remote"`
	Host   string `json:"host"`
	Owner  string `json:"owner"`
	Name   string `json:"name"`
	URL    string `json:"url"`
}
