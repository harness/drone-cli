package template

import (
	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/drone-go/drone"
	"github.com/urfave/cli"
	"io/ioutil"
	"strings"
)

var templateUpdateCmd = cli.Command{
	Name:      "update",
	Usage:     "update a template",
	ArgsUsage: "[name]",
	Action:    templateUpdate,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "template name",
		},
		cli.StringFlag{
			Name:  "data",
			Usage: "template file data",
		},
	},
}

func templateUpdate(c *cli.Context) error {
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}
	template := &drone.Template{
		Name: c.String("name"),
	}
	if strings.HasPrefix(c.String("data"), "@") {
		path := strings.TrimPrefix(c.String("data"), "@")
		out, ferr := ioutil.ReadFile(path)
		if ferr != nil {
			return ferr
		}
		template.Data = out
	}
	_, err = client.TemplateUpdate(template.Name, template)
	return err
}
