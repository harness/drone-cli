package common

type Config struct {
	Clone Step
	Build Step

	Compose map[string]Step
	Publish map[string]Step
	Deploy  map[string]Step
	Notify  map[string]Step

	Matrix map[string][]string
}

type Step struct {
	Name        string
	Image       string
	Environment []string
	Volumes     []string
	Hostname    string
	Privileged  bool
	Net         string

	// Config represents the unique configuration details
	// for each plugin.
	Config map[string]interface{} `yaml:"config,inline"`
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
