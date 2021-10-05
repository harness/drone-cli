package build

import (
	"os"
	"text/template"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/funcmap"
	"github.com/urfave/cli"
)

var buildQueueV2Cmd = cli.Command{
	Name:      "queue-v2",
	Usage:     "show build queue",
	ArgsUsage: "",
	Action:    buildQueueV2,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Usage: "format output",
			Value: tmplQueueV2Status,
		},
		cli.StringFlag{
			Name:  "repo",
			Usage: "repo filter",
		},
	},
}

func buildQueueV2(c *cli.Context) error {
	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}

	instances, err := client.IncompleteV2()
	if err != nil {
		return err
	}

	tmpl, err := template.New("_").Funcs(funcmap.Funcs).Parse(c.String("format") + "\n")
	if err != nil {
		return err
	}

	slug := c.String("repo")

	for _, instance := range instances {
		if slug != "" && instance.RepoSlug != slug {
			continue
		}
		templateErr := tmpl.Execute(os.Stdout, instance)
		if templateErr != nil {
			return templateErr
		}
	}
	return nil
}

// template for build queue v2 information
var tmplQueueV2Status = "\x1b[33m{{ .RepoSlug }}#{{ .BuildNumber }} \x1b[0m" + `
Repo Namespace: {{ .RepoNamespace }}
Repo Name: {{ .RepoName }}
Repo Slug: {{ .RepoSlug }}
Build Number: {{ .BuildNumber }}
Build Author: {{ .BuildAuthor }}
Build Author Name : {{ .BuildAuthorName }}
Build Author Email: {{ .BuildAuthorEmail }}
Build Author Avatar : {{ .BuildAuthorAvatar }}
Build Sender: {{ .BuildSender }}
Build Started: {{ .BuildStarted | time }}
Build Finished: {{ .BuildFinished | time}}
Build Created: {{ .BuildCreated | time}}
Build Updated: {{ .BuildUpdated | time }}
Stage Name: {{ .StageName }}
Stage Kind: {{ .StageKind }}
Stage Type: {{ .StageType }}
Stage Status: {{ .StageStatus }}
Stage Machine: {{ .StageMachine }}
Stage OS: {{ .StageOS }}
Stage Arch: {{ .StageArch }}
Stage Variant: {{ .StageVariant }}
Stage Kernel: {{ .StageKernel }}
Stage Limit: {{ .StageLimit }}
Stage Limit Repo: {{ .StageLimitRepo }}
Stage Started: {{ .StageStarted | time }}
Stage Stopped: {{ .StageStopped | time }}
`
