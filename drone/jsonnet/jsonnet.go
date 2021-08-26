package jsonnet

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/urfave/cli"
)

// Command exports the jsonnet command.
var Command = cli.Command{
	Name:      "jsonnet",
	Usage:     "generate .drone.yml from jsonnet",
	ArgsUsage: "[path/to/.drone.jsonnet]",
	Action: func(c *cli.Context) {
		if err := generate(c); err != nil {
			log.Fatalln(err)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "source",
			Usage: "Source file",
			Value: ".drone.jsonnet",
		},
		cli.StringFlag{
			Name:  "target",
			Usage: "target file",
			Value: ".drone.yml",
		},
		cli.BoolFlag{
			Name:  "stream",
			Usage: "Write output as a YAML stream.",
		},
		cli.BoolFlag{
			Name:  "format",
			Usage: "Write output as formatted YAML",
		},
		cli.BoolFlag{
			Name:  "stdout",
			Usage: "Write output to stdout",
		},
		cli.BoolFlag{
			Name:  "string",
			Usage: "Expect a string, manifest as plain text",
		},
		cli.StringSliceFlag{
			Name:  "extVar, V",
			Usage: "Pass extVars to Jsonnet (can be specified multiple times)",
		},
	},
}

func generate(c *cli.Context) error {
	result, err := convert(c.String("source"), c.Bool("string"), c.Bool("format"), c.Bool("stream"), c.StringSlice("extVar"))
	if err != nil {
		return err
	}

	// the user can optionally write the yaml to stdout. This is useful for debugging purposes without mutating an existing file.
	if c.Bool("stdout") {
		io.WriteString(os.Stdout, result)
		return nil
	}

	target := c.String("target")
	return ioutil.WriteFile(target, []byte(result), 0644)
}
