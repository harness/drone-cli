package sign

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/drone/drone-cli/drone/internal"
	"github.com/drone/drone-yaml/yaml"
	"github.com/drone/drone-yaml/yaml/pretty"
	"github.com/urfave/cli"
)

// Command exports the sign command.
var Command = cli.Command{
	Name:      "sign",
	Usage:     "sign the yaml file",
	ArgsUsage: "<source>",
	Action:    format,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "save",
			Usage: "save result to source",
		},
	},
}

func format(c *cli.Context) error {
	repo := c.Args().First()
	owner, name, err := internal.ParseRepo(repo)
	if err != nil {
		return err
	}

	client, err := internal.NewClient(c)
	if err != nil {
		return err
	}

	path := c.Args().Get(1)
	if path == "" {
		path = ".drone.yml"
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	hmac, err := client.Sign(owner, name, string(data))
	if err != nil {
		return err
	}

	if c.Bool("save") == false {
		fmt.Println(hmac)
		return nil
	}

	manifest, err := yaml.ParseFile(path)
	if err != nil {
		return err
	}

	var append bool
	for _, resource := range manifest.Resources {
		signature, ok := resource.(*yaml.Signature)
		if ok {
			append = false
			signature.Hmac = hmac
			break
		}
	}
	if append {
		// TODO this is currently disabled becuase it is resulting
		// in a compiler panic. I need to investigate this further.

		// manifest.Resources = append(
		// 	manifest.Resources,
		// 	&yaml.Signature{
		// 		Kind: "signature",
		// 		Hmac: hmac,
		// 	},
		// )
	}

	buf := new(bytes.Buffer)
	pretty.Print(buf, manifest)
	return ioutil.WriteFile(path, buf.Bytes(), 0644)
}
