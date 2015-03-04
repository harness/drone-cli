package common

type Config struct {
	/*
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
	*/

	//Clone   map[string]Step
	//Build   map[string]Step

	Clone Step
	Build Step

	Compose map[string]Step
	Publish map[string]Step //interface{}
	Deploy  map[string]Step //interface{}
	Notify  map[string]Step //interface{}

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

	// Config represents the unique configuration details
	// for each plugin.
	Config map[string]interface{} `yaml:"config,inline" json:"config"`
}

/*

type Container interface {
	Create()
	Start()
	Wait()
	Logs()
	Info()
	Stop()
	Kill()
	Remove()
}

clone:
  git:
	  path: foo/bar

build:
  script:
	  image: go:1.2
	  commands:
		  - go build
			- go test

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
