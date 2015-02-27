package common

type User struct {
	Remote string `json:"remote"`
	Login  string `json:"login"`
	Name   string `json:"name"`
	Email  string `json:"email,omitempty"`
}
