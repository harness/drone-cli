package jsonnet

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"fmt"
	"regexp"

	"github.com/drone/drone-cli/drone/internal"
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
	ArgsUsage: "[path/to/repo]",
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
		cli.BoolFlag{
			Name:  "sign",
			Usage: "Sign jsonnet file and save",
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

	// the user can optionally sign and write the hmac
	// generated from the signing process to the jsonnet
	// file in a way that compiles to yaml as expected
	if c.Bool("sign") {
		repo := c.Args().First()
		owner, name, err := internal.ParseRepo(repo)
		if err != nil {
			return err
		}
		client, err := internal.NewClient(c)
		if err != nil {
			return err
		}
		hmac, err := client.Sign(owner, name, string(buf.Bytes()))
		if err != nil {
			return err
		}
		jsonnetBuf := new(bytes.Buffer)
		r, _ := regexp.Compile(`\s*\+\s*\[\s*\{\s*kind\s*:\s*"signature"\s*,(.|\s)*}\s*\]`)
		jsonnetBuf.Write(r.ReplaceAll(data, []byte("")))
		jsonnetBuf.WriteString(fmt.Sprintf(" + [{kind: \"signature\",hmac: \"%s\"}]\n", hmac))
		return ioutil.WriteFile(source, jsonnetBuf.Bytes(), 0644)
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
