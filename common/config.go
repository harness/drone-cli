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

type Step struct {
	Name        string   `json:"name"`
	Image       string   `json:"image"`
	Environment []string `json:"environment"`
	Volumes     []string `json:"volumes"`
	Hostname    string   `json:"hostname"`
	Privileged  bool     `json:"privileged"`
	Net         string   `json:"net"`

	Config map[string]interface{} `yaml:"config,inner" json:"config"`
}

/*

type Step struct {
	Image       string
	Commands    []string
	Links       []string
	Ports       []string
	Expose      []string
	Volumes     []string
	VolumesFrom []string
	Environment []string
	DNS         []string
	Net         string

	CapAdd     []string
	CapDrop    []string
	Workdir    string
	Entrypoint string
	User       string
	Hostname   string
	Privileged bool

	MemLimit  int64
	CPUShared int
}
*/
