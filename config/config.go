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
		Privileged bool
		Hostname   string
		Net        string
	}

	Notify  map[string]interface{}
	Publish map[string]interface{}
	Deploy  map[string]interface{}

	Matrix map[string][]string
}

// func (c *Config) GetEnv() string    { return c.Env }
// func (c *Config) GetImage() string  { return c.Image }
// func (c *Config) GetScript() string { return c.Script }

//func (c *Config) GetDocker() {
//	return
//}
