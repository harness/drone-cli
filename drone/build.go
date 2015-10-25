package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/drone/drone-go/drone"
)

var BuildCmd = cli.Command{
	Name:  "build",
	Usage: "manage builds",
	Subcommands: []cli.Command{
		// build info
		{
			Name:  "info",
			Usage: "show build information",
			Action: func(c *cli.Context) {
				handle(c, BuildInfoCmd)
			},
		},
		// build list
		{
			Name:  "list",
			Usage: "list recent builds",
			Action: func(c *cli.Context) {
				handle(c, BuildListCmd)
			},
		},
		// build logs
		{
			Name:  "logs",
			Usage: "show build logs",
			Action: func(c *cli.Context) {
				handle(c, BuildLogCmd)
			},
		},
		// build start
		{
			Name:  "start",
			Usage: "start a stopped build",
			Action: func(c *cli.Context) {
				handle(c, BuildStartCmd)
			},
		},
		// build stop
		{
			Name:  "stop",
			Usage: "stop a running build job",
			Action: func(c *cli.Context) {
				handle(c, BuildStopCmd)
			},
		},
	},
}

func BuildInfoCmd(c *cli.Context, client drone.Client) error {
	var (
		nameParam = c.Args().Get(0)
		numParam  = c.Args().Get(1)

		err   error
		owner string
		name  string
		num   int
	)

	num, err = strconv.Atoi(numParam)
	if err != nil {
		return fmt.Errorf("Invalid or missing build number")
	}

	owner, name, err = parseRepo(nameParam)
	if err != nil {
		return err
	}
	build, err := client.Build(owner, name, num)
	if err != nil {
		return err
	}
	fmt.Println(build.Number)
	fmt.Println(build.Event)
	fmt.Println(build.Status)
	fmt.Println(build.Created)
	fmt.Println(build.Started)
	fmt.Println(build.Enqueued)
	fmt.Println(build.Finished)
	fmt.Println(build.Commit)
	fmt.Println(build.Ref)
	fmt.Println(build.Author)
	fmt.Println(build.Message)
	return nil
}

func BuildListCmd(c *cli.Context, client drone.Client) error {
	owner, name, err := parseRepo(c.Args().Get(0))
	if err != nil {
		return err
	}
	builds, err := client.BuildList(owner, name)
	if err != nil {
		return err
	}

	for _, build := range builds {
		fmt.Println(build.Number)
	}
	return nil
}

func BuildLogCmd(c *cli.Context, client drone.Client) error {
	var (
		nameParam = c.Args().Get(0)
		numParam  = c.Args().Get(1)
		jobParam  = c.Args().Get(2)

		err   error
		owner string
		name  string
		num   int
		job   int
	)

	num, err = strconv.Atoi(numParam)
	if err != nil {
		return fmt.Errorf("Invalid or missing build number")
	}
	job, err = strconv.Atoi(jobParam)
	if err != nil {
		return fmt.Errorf("Invalid or missing job number")
	}

	owner, name, err = parseRepo(nameParam)
	if err != nil {
		return err
	}
	rc, err := client.BuildLogs(owner, name, num, job)
	if err != nil {
		return err
	}
	defer rc.Close()

	io.Copy(os.Stdout, rc)
	return nil
}

func BuildStartCmd(c *cli.Context, client drone.Client) error {
	var (
		nameParam = c.Args().Get(0)
		numParam  = c.Args().Get(1)

		err   error
		owner string
		name  string
		num   int
	)

	num, err = strconv.Atoi(numParam)
	if err != nil {
		return fmt.Errorf("Invalid or missing build number")
	}
	owner, name, err = parseRepo(nameParam)
	if err != nil {
		return err
	}

	build, err := client.BuildStart(owner, name, num)
	if err != nil {
		return err
	}

	fmt.Println(build.Number)
	fmt.Println(build.Status)

	return nil
}

func BuildStopCmd(c *cli.Context, client drone.Client) error {
	var (
		nameParam = c.Args().Get(0)
		numParam  = c.Args().Get(1)
		jobParam  = c.Args().Get(2)

		err   error
		owner string
		name  string
		num   int
		job   int
	)

	// jobs are really only useful to matrix builds. So if the
	// job is not specified we'll assume it isn't a matrix build
	// and we'll cancel job number 1
	if len(jobParam) == 0 {
		jobParam = "1"
	}

	num, err = strconv.Atoi(numParam)
	if err != nil {
		return fmt.Errorf("Invalid or missing build number")
	}
	job, err = strconv.Atoi(jobParam)
	if err != nil {
		return fmt.Errorf("Invalid or missing job number")
	}
	owner, name, err = parseRepo(nameParam)
	if err != nil {
		return err
	}

	err = client.BuildStop(owner, name, num, job)
	if err != nil {
		return err
	}

	fmt.Printf("stopping build %d job %d\n", num, job)
	return nil
}
