package jsonnet

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/drone/drone-cli/drone/jsonnet/stdlib"
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
			Name: "stream",
		},
		cli.BoolTFlag{
			Name: "string",
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
	vm.MaxStack = c.Int("max-stack")
	vm.StringOutput = c.BoolT("string")
	vm.ErrorFormatter.SetMaxStackTraceSize(
		c.Int("max-trace"),
	)
	vm.ErrorFormatter.SetColorFormatter(
		color.New(color.FgRed).Fprintf,
	)
	vm.Importer(
		stdlib.Importer(),
	)

	var result string
	var resultArray []string
	if c.Bool("stream") {
		resultArray, err = vm.EvaluateSnippetStream(input, string(snippet))
	} else {
		result, err = vm.EvaluateSnippet(input, string(snippet))
	}
	if err != nil {
		return err
	}

	if c.Bool("stream") {
		return writeOutputStream(resultArray, output)
	}
	return writeOutputFile(result, output)
}

// writeOutputStream writes the output as a YAML stream.
func writeOutputStream(output []string, outputFile string) error {
	var f *os.File

	if outputFile == "" {
		f = os.Stdout
	} else {
		var err error
		f, err = os.Create(outputFile)
		if err != nil {
			return err
		}
		defer f.Close()
	}

	for _, doc := range output {
		_, err := f.WriteString("---\n")
		if err != nil {
			return err
		}
		_, err = f.WriteString(doc)
		if err != nil {
			return err
		}
	}

	if len(output) > 0 {
		_, err := f.WriteString("...\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func writeOutputFile(output string, outputFile string) error {
	if outputFile == "" {
		fmt.Print(output)
		return nil
	}

	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(output)
	return err
}
