package jsonnet

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/ghodss/yaml"
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
		cli.StringSliceFlag{
			Name:  "jpath, J",
			Usage: "Specify an additional library search dir (right-most wins)",
		},
	},
}

func generate(c *cli.Context) error {
	result, err := convert(
		c.String("source"),
		c.Bool("string"),
		c.Bool("format"),
		c.Bool("stream"),
		c.StringSlice("extVar"),
		c.StringSlice("jpath"),
	)
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

func convert(source string, stringOutput bool, format bool, stream bool, vars []string, jpath []string) (string, error) {
	data, err := ioutil.ReadFile(source)
	if err != nil {
		return "", err
	}

	vm := jsonnet.MakeVM()
	vm.MaxStack = 500
	vm.StringOutput = stringOutput
	vm.ErrorFormatter.SetMaxStackTraceSize(20)
	vm.ErrorFormatter.SetColorFormatter(
		color.New(color.FgRed).Fprintf,
	)

	// register native functions
	RegisterNativeFuncs(vm)

	jsonnetPath := filepath.SplitList(os.Getenv("JSONNET_PATH"))
	jsonnetPath = append(jsonnetPath, jpath...)
	vm.Importer(&jsonnet.FileImporter{
		JPaths: jsonnetPath,
	})

	// extVars
	for _, v := range vars {
		name, value, err := getVarVal(v)
		if err != nil {
			return "", err
		}
		vm.ExtVar(name, value)
	}

	formatDoc := func(doc []byte) ([]byte, error) {
		// enable yaml output
		if format {
			formatted, yErr := yaml.JSONToYAML(doc)
			if yErr != nil {
				return nil, fmt.Errorf("failed to convert to YAML: %v", yErr)
			}
			return formatted, nil
		}
		return doc, nil
	}

	buf := new(bytes.Buffer)
	if stream {
		docs, err := vm.EvaluateSnippetStream(source, string(data))
		if err != nil {
			return "", err
		}
		for _, doc := range docs {
			formatted, err := formatDoc([]byte(doc))
			if err != nil {
				return "", err
			}

			buf.WriteString("---")
			buf.WriteString("\n")
			buf.Write(formatted)
		}
	} else {
		result, err := vm.EvaluateSnippet(source, string(data))
		if err != nil {
			return "", err
		}
		formatted, err := formatDoc([]byte(result))
		if err != nil {
			return "", err
		}
		buf.Write(formatted)
	}

	return buf.String(), nil
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
