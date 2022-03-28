package template

import (
	"errors"
	"fmt"
	"github.com/drone/envsubst"
	"github.com/drone/runner-go/manifest"
	"html/template"
	"io/ioutil"
	"strings"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/funcmap"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"os"
)

var templateRenderCmd = cli.Command{
	Name:      "render",
	Usage:     "render a pipeline from a template",
	ArgsUsage: "[namespace] [name] [input]",
	Action:    templateRender,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "namespace",
			Usage: "organization namespace",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "template name",
		},
		cli.StringFlag{
			Name:  "input",
			Usage: "path to input file",
		},
	},
}

type templateArgs struct {
	Kind string
	Load string
	Data map[string]interface{}
}

/*
Example usage:

drone template render --namespace foo --name my_template.yaml --input my_values.yaml

Where "my_values.yaml" is a valid YAML file containing a `data` key:

```
data:
  foo: bar
```

*/
func templateRender(c *cli.Context) error {
	var (
		namespace    = c.String("namespace")
		templateName = c.String("name")
		input        = c.String("input")
	)
	if templateName == "" {
		return errors.New("Missing template name")
	}
	if namespace == "" {
		return errors.New("Missing namespace")
	}

	var envs map[string]string

	// Find pipeline
	rawSource, err := ioutil.ReadFile(input)
	if err != nil {
		return err
	}

	// string substitution function ensures that string replacement variables are escaped and quoted if they contain newlines.
	subf := func(k string) string {
		v := envs[k]
		if strings.Contains(v, "\n") {
			v = fmt.Sprintf("%q", v)
		}
		return v
	}
	// evaluates string replacement expressions and returns an update configuration.
	_, err = envsubst.Eval(string(rawSource), subf)
	if err != nil {
		return err
	}

	// we need to parse the file again into raw resources to access the dependencies
	inputRawResources, err := manifest.ParseRawFile(input)
	if err != nil {
		return err
	}

	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}
	templates, err := client.Template(namespace, templateName)
	if err != nil {
		return err
	}
	tmpl, err := template.New("_").Funcs(funcmap.Funcs).Parse(templates.Data)
	if err != nil {
		return err
	}

	var inputArgs templateArgs

	err = yaml.Unmarshal([]byte(inputRawResources[0].Data), &inputArgs)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"input": inputArgs.Data,
	}

	return tmpl.Execute(os.Stdout, data)
}
