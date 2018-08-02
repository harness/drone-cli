package jsonnet

import (
	"io/ioutil"
	"log"

	"github.com/fatih/color"
	"github.com/google/go-jsonnet"
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
		cli.BoolFlag{
			Name:   "string",
			Hidden: true,
		},
		cli.IntFlag{
			Name:  "max-stack",
			Usage: "number of allowed stack frames",
			Value: 500,
		},
		cli.IntFlag{
			Name:  "max-trace",
			Usage: "max length of stack trace before cropping",
			Value: 20,
		},
		cli.BoolFlag{
			Name:  "stdout",
			Usage: "write the json document to stdout",
		},
	},
}

func generate(c *cli.Context) error {
	input := c.Args().Get(0)
	if input == "" {
		input = ".drone.jsonnet"
	}
	output := c.Args().Get(1)
	if output == "" {
		output = ".drone.yml"
	}

	snippet, err := ioutil.ReadFile(input)
	if err != nil {
		return err
	}

	vm := jsonnet.MakeVM()
	vm.KeepOrder = true
	vm.StringOutput = c.Bool("string")
	vm.MaxStack = c.Int("max-stack")
	// vm.ExtVar
	// vm.ExtCode
	// vm.TLAVar
	// vm.TLACode
	// vm.Importer(&jsonnet.FileImporter{})

	vm.ErrorFormatter.SetMaxStackTraceSize(
		c.Int("max-trace"),
	)
	vm.ErrorFormatter.SetColorFormatter(
		color.New(color.FgRed).Fprintf,
	)

	raw, err := vm.EvaluateSnippet(input, string(snippet))
	if err != nil {
		return err
	}

	if c.Bool("stdout") {
		println(raw)
		return nil
	}
	return ioutil.WriteFile(output, []byte(raw), 0644)
}
