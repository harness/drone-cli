package format

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/urfave/cli"
)

// Command exports the fmt command.
var Command = cli.Command{
	Name:      "fmt",
	Usage:     "<deprecated. this operation is a no-op> format the yaml file",
	ArgsUsage: "<source>",
	Hidden:    true,
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
	out, _ := ioutil.ReadFile(path)
	buf := bytes.NewBuffer(out)
	if c.Bool("save") {
		return ioutil.WriteFile(path, buf.Bytes(), 0644)
	}
	_, err := io.Copy(os.Stderr, buf)
	return err
}
