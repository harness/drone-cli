package convert

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/drone/drone-yaml/yaml/converter"
	"github.com/urfave/cli"
)

// Command exports the convert command.
var Command = cli.Command{
	Name:      "convert",
	Usage:     "convert legacy format",
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

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	res, err := converter.Convert(raw, converter.Metadata{Filename: path})
	if err != nil {
		return err
	}

	if c.Bool("save") {
		return ioutil.WriteFile(path, res, 0644)
	}

	_, err = io.Copy(os.Stderr, bytes.NewReader(res))
	return err
}
