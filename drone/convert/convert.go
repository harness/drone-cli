package convert

import (
	"io/ioutil"

	"github.com/urfave/cli"
)

// Command exports the convert command.
var Command = cli.Command{
	Name:      "convert",
	Usage:     "<deprecated. this operation is a no-op> convert legacy format",
	ArgsUsage: "<source>",
	Action:    convert,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "save",
			Usage: "save result to source",
		},
	},
}

func convert(c *cli.Context) error {
	path := c.Args().First()
	if path == "" {
		path = ".drone.yml"
	}

	_, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return err
}
