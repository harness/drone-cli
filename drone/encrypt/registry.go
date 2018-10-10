package encrypt

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"strings"

// 	"github.com/drone/drone-cli/drone/internal"
// 	"github.com/drone/drone-go/drone"

// 	"github.com/urfave/cli"
// )

// var encryptRegistryCommand = cli.Command{
// 	Name:      "registry",
// 	Usage:     "encrypt registry credentials",
// 	ArgsUsage: "<repo/name> <string>",
// 	Action:    encryptRegistry,
// 	Flags: []cli.Flag{
// 		cli.StringFlag{
// 			Name:  "username",
// 			Usage: "registry username",
// 		},
// 		cli.StringFlag{
// 			Name:  "password",
// 			Usage: "registry password",
// 		},
// 		cli.StringFlag{
// 			Name:  "server",
// 			Usage: "registry server",
// 			Value: "docker.io",
// 		},
// 	},
// }

// func encryptRegistry(c *cli.Context) error {
// 	repo := c.Args().First()
// 	owner, name, err := internal.ParseRepo(repo)
// 	if err != nil {
// 		return err
// 	}

// 	client, err := internal.NewClient(c)
// 	if err != nil {
// 		return err
// 	}

// 	password := c.String("password")
// 	if strings.HasPrefix(password, "@") {
// 		data, err := ioutil.ReadFile(password)
// 		if err != nil {
// 			return err
// 		}
// 		password = string(data)
// 	}

// 	policy := "pull"
// 	switch {
// 	case c.Bool("push"):
// 		policy = "push"
// 	case c.Bool("push-pull-request"):
// 		policy = "push-pull-request"
// 	}

// 	registry := &drone.Registry{
// 		Address:  c.String("server"),
// 		Username: c.String("username"),
// 		Password: password,
// 		Policy:   policy,
// 	}
// 	encrypted, err := client.EncryptRegistry(owner, name, registry)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println(encrypted)
// 	return nil
// }
