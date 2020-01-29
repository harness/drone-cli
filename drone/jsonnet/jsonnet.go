package jsonnet

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/drone/drone-yaml/yaml"
	"github.com/drone/drone-yaml/yaml/pretty"
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
		cli.BoolTFlag{
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
	source := c.String("source")
	target := c.String("target")

	data, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}

	vm := jsonnet.MakeVM()
	vm.MaxStack = 500
	vm.StringOutput = c.Bool("string")
	vm.ErrorFormatter.SetMaxStackTraceSize(20)
	vm.ErrorFormatter.SetColorFormatter(
		color.New(color.FgRed).Fprintf,
	)

	// register native functions
	RegisterNativeFuncs(vm)

	// extVars
	vars := c.StringSlice("extVar")
	for _, v := range vars {
		name, value, err := getVarVal(v)
		if err != nil {
			return err
		}
		vm.ExtVar(name, value)
	}

	buf := new(bytes.Buffer)
	if c.Bool("stream") {
		docs, err := vm.EvaluateSnippetStream(source, string(data))
		if err != nil {
			return err
		}
		for _, doc := range docs {
			buf.WriteString("---")
			buf.WriteString("\n")
			buf.WriteString(doc)
		}
	} else {
		result, err := vm.EvaluateSnippet(source, string(data))
		if err != nil {
			return err
		}
		buf.WriteString(result)
	}

	// the yaml file is parsed and formatted by default. This
	// can be disabled for --format=false.
	if c.BoolT("format") {
		manifest, err := yaml.Parse(buf)
		if err != nil {
			return err
		}
		buf.Reset()
		pretty.Print(buf, manifest)
	}

	// the user can optionally write the yaml to stdout. This
	// is useful for debugging purposes without mutating an
	// existing file.
	if c.Bool("stdout") {
		io.Copy(os.Stdout, buf)
		return nil
	}

	return ioutil.WriteFile(target, buf.Bytes(), 0644)
}

// https://github.com/google/go-jsonnet/blob/master/cmd/jsonnet/cmd.go#L149
func getVarVal(s string) (string, string, error) {
	parts := strings.SplitN(s, "=", 2)
	name := parts[0]
	if len(parts) == 1 {
		content, exists := os.LookupEnv(name)
		if exists {
			return name, content, nil
		}
		return "", "", fmt.Errorf("environment variable %v was undefined", name)
	}
	return name, parts[1], nil
}
