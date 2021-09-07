package starlark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/ghodss/yaml"
	"github.com/urfave/cli"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

// Command exports the jsonnet command.
var Command = cli.Command{
	Name:      "starlark",
	Usage:     "generate .drone.yml from starlark",
	ArgsUsage: "[path/to/.drone.star]",
	Action: func(c *cli.Context) {
		if err := generate(c); err != nil {
			if err, ok := err.(*starlark.EvalError); ok {
				log.Fatalf("starlark evaluation error:\n%s", err.Backtrace())
			}
			log.Fatalln(err)
		}
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "source",
			Usage: "Source file",
			Value: ".drone.star",
		},
		cli.StringFlag{
			Name:  "target",
			Usage: "target file",
			Value: ".drone.yml",
		},
		cli.BoolFlag{
			Name:  "format",
			Usage: "Write output as formatted YAML",
		},
		cli.BoolFlag{
			Name:  "stdout",
			Usage: "Write output to stdout",
		},
		cli.Uint64Flag{
			Name:  "max-execution-steps",
			Usage: "maximum number of execution steps",
			Value: 50000,
		},
		//
		// Drone Parameters
		//
		cli.StringFlag{
			Name:  "repo.name",
			Usage: "repository name",
		},
		cli.StringFlag{
			Name:  "repo.namespace",
			Usage: "repository namespace",
		},
		cli.StringFlag{
			Name:  "repo.slug",
			Usage: "repository slug",
		},

		cli.StringFlag{
			Name:  "build.event",
			Usage: "build event",
			Value: "push",
		},
		cli.StringFlag{
			Name:  "build.branch",
			Usage: "build branch",
			Value: "master",
		},
		cli.StringFlag{
			Name:  "build.source",
			Usage: "build source branch",
			Value: "master",
		},
		cli.StringFlag{
			Name:  "build.source_repo",
			Usage: "repo slug of source repository",
		},
		cli.StringFlag{
			Name:  "build.target",
			Usage: "build target branch",
			Value: "master",
		},
		cli.StringFlag{
			Name:  "build.ref",
			Usage: "build ref",
			Value: "refs/heads/master",
		},
		cli.StringFlag{
			Name:  "build.commit",
			Usage: "build commit sha",
		},
		cli.StringFlag{
			Name:  "build.message",
			Usage: "build commit message",
		},
		cli.StringFlag{
			Name:  "build.title",
			Usage: "build title",
		},
		cli.StringFlag{
			Name:  "build.link",
			Usage: "build link",
		},
		cli.StringFlag{
			Name:  "build.environment",
			Usage: "build environment",
		},
		cli.BoolTFlag{
			Name:  "build.debug",
			Usage: "debug build",
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

	thread := &starlark.Thread{
		Name:  "drone",
		Print: func(_ *starlark.Thread, msg string) { fmt.Println(msg) },
		Load:  makeLoad(),
	}
	globals, err := starlark.ExecFile(thread, source, data, nil)
	if err != nil {
		return err
	}

	mainVal, ok := globals["main"]
	if !ok {
		return fmt.Errorf("no main function found")
	}
	main, ok := mainVal.(starlark.Callable)
	if !ok {
		return fmt.Errorf("main must be a function")
	}

	// TODO this needs to be flushed out.
	repo := starlark.StringDict{
		"name":      starlark.String(c.String("repo.name")),
		"namespace": starlark.String(c.String("repo.namespace")),
		"slug":      starlark.String(c.String("repo.slug")),
	}

	build := starlark.StringDict{
		"event":       starlark.String(c.String("build.event")),
		"branch":      starlark.String(c.String("build.branch")),
		"source":      starlark.String(c.String("build.source_branch")),
		"source_repo": starlark.String(c.String("build.source_repo")),
		"target":      starlark.String(c.String("build.target_branch")),
		"ref":         starlark.String(c.String("build.ref")),
		"commit":      starlark.String(c.String("build.commit")),
		"message":     starlark.String(c.String("build.message")),
		"title":       starlark.String(c.String("build.title")),
		"link":        starlark.String(c.String("build.link")),
		"environment": starlark.String(c.String("build.environment")),
		"debug":       starlark.Bool(c.Bool("build.debug")),
	}

	args := starlark.Tuple([]starlark.Value{
		starlarkstruct.FromStringDict(
			starlark.String("context"),
			starlark.StringDict{
				"repo":  starlarkstruct.FromStringDict(starlark.String("repo"), repo),
				"build": starlarkstruct.FromStringDict(starlark.String("build"), build),
			},
		),
	})
	thread.SetMaxExecutionSteps(c.Uint64("max-execution-steps"))
	mainVal, err = starlark.Call(thread, main, args, nil)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	switch v := mainVal.(type) {
	case *starlark.List:
		for i := 0; i < v.Len(); i++ {
			tmpBuf := new(bytes.Buffer)
			item := v.Index(i)
			tmpBuf.WriteString("---\n")
			err = writeJSON(tmpBuf, item)
			if err != nil {
				return err
			}
			tmpBuf.WriteString("\n")
			if c.Bool("format") {
				yml, yErr := yaml.JSONToYAML(tmpBuf.Bytes())
				if yErr != nil {
					return fmt.Errorf("failed to convert to YAML: %v", yErr)
				}
				tmpBuf.Reset()
				tmpBuf.WriteString("---\n")
				tmpBuf.Write(yml)
			}
			buf.Write(tmpBuf.Bytes())
		}
	case *starlark.Dict:
		buf.WriteString("---\n")
		err = writeJSON(buf, v)
		if err != nil {
			return err
		}
		if c.BoolT("format") {
			yml, yErr := yaml.JSONToYAML(buf.Bytes())
			if yErr != nil {
				return fmt.Errorf("failed to convert to YAML: %v", yErr)
			}
			buf.Reset()
			buf.Write(yml)
		}
	default:
		return fmt.Errorf("invalid return type (got a %s)", mainVal.Type())
	}
	// the user can optionally write the to stdout. This is useful for debugging purposes without mutating an existing file.
	if c.Bool("stdout") {
		io.Copy(os.Stdout, buf)
		return nil
	}
	return ioutil.WriteFile(target, buf.Bytes(), 0644)
}

// Adapted from skycfg:
// https://github.com/stripe/skycfg/blob/eaa524101c2a0807c13ed5d2e52576fefc146ec3/internal/go/skycfg/json_write.go#L45
func writeJSON(out *bytes.Buffer, v starlark.Value) error {
	if marshaler, ok := v.(json.Marshaler); ok {
		jsonData, err := marshaler.MarshalJSON()
		if err != nil {
			return err
		}
		out.Write(jsonData)
		return nil
	}

	switch v := v.(type) {
	case starlark.NoneType:
		out.WriteString("null")
	case starlark.Bool:
		fmt.Fprintf(out, "%t", v)
	case starlark.Int:
		out.WriteString(v.String())
	case starlark.Float:
		fmt.Fprintf(out, "%g", v)
	case starlark.String:
		s := string(v)
		if goQuoteIsSafe(s) {
			fmt.Fprintf(out, "%q", s)
		} else {
			// vanishingly rare for text strings
			data, _ := json.Marshal(s)
			out.Write(data)
		}
	case starlark.Indexable: // Tuple, List
		out.WriteByte('[')
		for i, n := 0, starlark.Len(v); i < n; i++ {
			if i > 0 {
				out.WriteString(", ")
			}
			if err := writeJSON(out, v.Index(i)); err != nil {
				return err
			}
		}
		out.WriteByte(']')
	case *starlark.Dict:
		out.WriteByte('{')
		for i, itemPair := range v.Items() {
			key := itemPair[0]
			value := itemPair[1]
			if i > 0 {
				out.WriteString(", ")
			}
			if err := writeJSON(out, key); err != nil {
				return err
			}
			out.WriteString(": ")
			if err := writeJSON(out, value); err != nil {
				return err
			}
		}
		out.WriteByte('}')
	default:
		return fmt.Errorf("TypeError: value %s (type `%s') can't be converted to JSON.", v.String(), v.Type())
	}
	return nil
}

func goQuoteIsSafe(s string) bool {
	for _, r := range s {
		// JSON doesn't like Go's \xHH escapes for ASCII control codes,
		// nor its \UHHHHHHHH escapes for runes >16 bits.
		if r < 0x20 || r >= 0x10000 {
			return false
		}
	}
	return true
}

// https://github.com/google/starlark-go/blob/4eb76950c5f02ec5bcfd3ca898231a6543942fd9/repl/repl.go#L175
func makeLoad() func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	type entry struct {
		globals starlark.StringDict
		err     error
	}

	var cache = make(map[string]*entry)

	return func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
		e, ok := cache[module]
		if e == nil {
			if ok {
				// request for package whose loading is in progress
				return nil, fmt.Errorf("cycle in load graph")
			}

			// Add a placeholder to indicate "load in progress".
			cache[module] = nil

			// Load it.
			thread := &starlark.Thread{Name: "exec " + module, Load: thread.Load}
			globals, err := starlark.ExecFile(thread, module, nil, nil)
			e = &entry{globals, err}

			// Update the cache.
			cache[module] = e
		}
		return e.globals, e.err
	}
}
