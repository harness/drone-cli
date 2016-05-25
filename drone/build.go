package main

import (
	"fmt"
	"html/template"
	"os"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/drone/drone-go/drone"
)

var buildCmd = cli.Command{
	Name:  "build",
	Usage: "manage builds",
	Subcommands: []cli.Command{

		// list command
		cli.Command{
			Name:   "list",
			Usage:  "show build history",
			Action: buildList,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "format",
					Usage: "format output",
					Value: tmplBuildList,
				},
				cli.StringFlag{
					Name:  "branch",
					Usage: "branch filter",
				},
				cli.StringFlag{
					Name:  "event",
					Usage: "event filter",
				},
				cli.StringFlag{
					Name:  "status",
					Usage: "status filter",
				},
				cli.IntFlag{
					Name:  "limit",
					Usage: "limit the list size",
					Value: 25,
				},
			},
		},

		// last command
		cli.Command{
			Name:   "last",
			Usage:  "show latest build details",
			Action: buildLast,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "format",
					Usage: "format output",
					Value: tmplBuildInfo,
				},
				cli.StringFlag{
					Name:  "branch",
					Usage: "branch name",
					Value: "master",
				},
			},
		},

		// info command
		cli.Command{
			Name:   "info",
			Usage:  "show build details",
			Action: buildInfo,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "format",
					Usage: "format output",
					Value: tmplBuildInfo,
				},
			},
		},

		// stop command
		cli.Command{
			Name:   "stop",
			Usage:  "stop a build",
			Action: buildStop,
		},

		// start command
		cli.Command{
			Name:   "start",
			Usage:  "start a build",
			Action: buildStart,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "fork",
					Usage: "fork the build",
				},
			},
		},

		// queue command.
		cli.Command{
			Name:   "queue",
			Usage:  "show build queue",
			Action: buildQueue,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "format",
					Usage: "format output",
					Value: tmplBuildQueue,
				},
			},
		},
	},
}

// command to display a list of recent builds.
func buildList(c *cli.Context) error {
	repo := c.Args().First()
	owner, name, err := parseRepo(repo)
	if err != nil {
		return err
	}

	client, err := newClient(c)
	if err != nil {
		return err
	}

	builds, err := client.BuildList(owner, name)
	if err != nil {
		return err
	}

	tmpl, err := template.New("_").Parse(c.String("format") + "\n")
	if err != nil {
		return err
	}

	branch := c.String("branch")
	event := c.String("event")
	status := c.String("status")
	limit := c.Int("limit")

	var count int
	for _, build := range builds {
		if count >= limit {
			break
		}
		if branch != "" && build.Branch != branch {
			continue
		}
		if event != "" && build.Event != event {
			continue
		}
		if status != "" && build.Status != status {
			continue
		}
		tmpl.Execute(os.Stdout, build)
		count++
	}
	return nil
}

// command to display the last build
func buildLast(c *cli.Context) error {
	repo := c.Args().First()
	owner, name, err := parseRepo(repo)
	if err != nil {
		return err
	}

	client, err := newClient(c)
	if err != nil {
		return err
	}

	build, err := client.BuildLast(owner, name, c.String("branch"))
	if err != nil {
		return err
	}

	tmpl, err := template.New("_").Parse(c.String("format"))
	if err != nil {
		return err
	}
	return tmpl.Execute(os.Stdout, build)
}

// command to display build information.
func buildInfo(c *cli.Context) error {
	repo := c.Args().First()
	owner, name, err := parseRepo(repo)
	if err != nil {
		return err
	}
	number, err := strconv.Atoi(c.Args().Get(1))
	if err != nil {
		return err
	}

	client, err := newClient(c)
	if err != nil {
		return err
	}

	build, err := client.Build(owner, name, number)
	if err != nil {
		return err
	}

	tmpl, err := template.New("_").Parse(c.String("format"))
	if err != nil {
		return err
	}
	return tmpl.Execute(os.Stdout, build)
}

//  command to stop a running build.
func buildStop(c *cli.Context) (err error) {
	repo := c.Args().First()
	owner, name, err := parseRepo(repo)
	if err != nil {
		return err
	}
	number, err := strconv.Atoi(c.Args().Get(1))
	if err != nil {
		return err
	}
	job, _ := strconv.Atoi(c.Args().Get(2))
	if job == 0 {
		job = 1
	}

	client, err := newClient(c)
	if err != nil {
		return err
	}

	err = client.BuildStop(owner, name, number, job)
	if err != nil {
		return err
	}

	fmt.Printf("Stopping build %s/%s#%d.%d\n", owner, name, number, job)
	return nil
}

// command to re-start an existing build.
func buildStart(c *cli.Context) (err error) {
	repo := c.Args().First()
	owner, name, err := parseRepo(repo)
	if err != nil {
		return err
	}
	number, err := strconv.Atoi(c.Args().Get(1))
	if err != nil {
		return err
	}

	client, err := newClient(c)
	if err != nil {
		return err
	}

	var build *drone.Build
	if c.Bool("fork") {
		build, err = client.BuildStart(owner, name, number)
	} else {
		build, err = client.BuildFork(owner, name, number)
	}
	if err != nil {
		return err
	}

	fmt.Printf("Starting build %s/%s#%d\n", owner, name, build.Number)
	return nil
}

// command to display the build queue.
func buildQueue(c *cli.Context) error {

	client, err := newClient(c)
	if err != nil {
		return err
	}

	builds, err := client.BuildQueue()
	if err != nil {
		return err
	}

	if len(builds) == 0 {
		fmt.Println("there are no pending or running builds")
		return nil
	}

	tmpl, err := template.New("_").Parse(c.String("format") + "\n")
	if err != nil {
		return err
	}

	for _, build := range builds {
		tmpl.Execute(os.Stdout, build)
	}
	return nil
}

// build info template.
var tmplBuildInfo = `Number: {{ .Number }}
Status: {{ .Status }}
Event: {{ .Event }}
Commit: {{ .Commit }}
Branch: {{ .Branch }}
Ref: {{ .Ref }}
Message: {{ .Message }}
Author: {{ .Author }}
`

// build queue template.
var tmplBuildQueue = "\x1b[33m{{ .FullName }} #{{ .Number }} \x1b[0m" + `
Status: {{ .Status }}
Event: {{ .Event }}
Commit: {{ .Commit }}
Branch: {{ .Branch }}
Ref: {{ .Ref }}
Author: {{ .Author }} {{ if .Email }}<{{.Email}}>{{ end }}
Message: {{ .Message }}
`

// build list template.
var tmplBuildList = "\x1b[33mBuild #{{ .Number }} \x1b[0m" + `
Status: {{ .Status }}
Event: {{ .Event }}
Commit: {{ .Commit }}
Branch: {{ .Branch }}
Ref: {{ .Ref }}
Author: {{ .Author }} {{ if .Email }}<{{.Email}}>{{ end }}
Message: {{ .Message }}
`
