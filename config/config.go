package config

type Config struct {
	Image    string
	Env      []string
	Script   []string
	Branches []string
	Services []string

	Git struct {
		Path string
	}

	Docker struct {
		Volumes    []string
		Privileged string
		Hostname   string
		Net        string
	}

	Notify  map[string]interface{}
	Publish map[string]interface{}
	Deploy  map[string]interface{}

	Matrix map[string][]string
}
