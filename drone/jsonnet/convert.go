package jsonnet

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/ghodss/yaml"
	"github.com/google/go-jsonnet"
)

func convert(source string, stringOutput bool, format bool, stream bool, vars []string) (string, error) {
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
