package format

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/drone/drone-yaml/yaml"
	"github.com/drone/drone-yaml/yaml/pretty"
	"github.com/urfave/cli"
)

// Command exports the fmt command.
var Command = cli.Command{
	Name:      "fmt",
	Usage:     "format the yaml file",
	ArgsUsage: "<source>",
	Action:    format,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "save",
			Usage: "save result to source",
		},
	},
}

func format(c *cli.Context) error {
	path := c.Args().First()
	if path == "" {
		path = ".drone.yml"
	}

	manifest, err := yaml.ParseFile(path)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	pretty.Print(buf, manifest)

	if c.Bool("save") {
		return ioutil.WriteFile(path, buf.Bytes(), 0644)
	}
	_, err = io.Copy(os.Stderr, buf)
	return err
}
