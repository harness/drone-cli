package template

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/drone-go/drone"
	"github.com/urfave/cli"
)

var templateCreateCmd = cli.Command{
	Name:      "add",
	Usage:     "adds a template",
	ArgsUsage: "[namespace] [name] [data]",
	Action:    templateCreate,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "template name",
		},
		cli.StringFlag{
			Name:  "namespace",
			Usage: "organization namespace",
		},
		cli.StringFlag{
			Name:  "data",
			Usage: "template file data",
		},
	},
}

func templateCreate(c *cli.Context) error {
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}
	template := &drone.Template{
		Namespace: c.String("namespace"),
		Name:      c.String("name"),
	}

	if template.Name == "" {
		return errors.New("Missing template name")
	}
	if template.Namespace == "" {
		return errors.New("Missing namespace")
	}

	if strings.HasPrefix(c.String("data"), "@") {
		path := strings.TrimPrefix(c.String("data"), "@")
		out, ferr := ioutil.ReadFile(path)
		if ferr != nil {
			return ferr
		}
		template.Data = string(out)
	}
	_, err = client.TemplateCreate(template)
	return err
}
