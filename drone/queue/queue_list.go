package queue

import (
	"fmt"
	"os"
	"text/template"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/funcmap"
	"github.com/urfave/cli"
)

var queueListCmd = cli.Command{
	Name:   "ls",
	Usage:  "list queue items",
	Action: queueList,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Usage: "format output",
			Value: tmplStage,
		},
	},
}

func queueList(c *cli.Context) (err error) {
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}
	builds, err := client.Queue()
	if err != nil {
		return err
	}

	if len(builds) == 0 {
		fmt.Println("there are no pending or running builds")
		return nil
	}

	tmpl, err := template.New("_").Funcs(funcmap.Funcs).Parse(c.String("format") + "\n")
	if err != nil {
		return err
	}

	for _, build := range builds {
		tmpl.Execute(os.Stdout, build)
	}
	return nil
}

var tmplStage = "\x1b[33mitem #{{ .ID }} \x1b[0m" + `
Status: {{ .Status }}
Machine: {{ .Machine }}
OS: {{ .OS }}
Arch: {{ .Arch }}
Variant: {{ .Variant }}
Version: {{ .Kernel }}
`
