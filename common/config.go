package common

type Config struct {
	Image    string
	Env      []string `json:"env"`
	Script   []string `json:"script"`
	Branches []string
	Services []string

	Git struct {
		Path string
	}

	Docker struct {
		Volumes    []string
		Privileged bool
		Hostname   string
		Net        string
	}

	Notify  map[string]interface{}
	Publish map[string]interface{}
	Deploy  map[string]interface{}

	Matrix map[string][]string
}
