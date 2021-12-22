package sign

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/buildkite/yaml"
	"github.com/drone/drone-cli/drone/internal"
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

	data, err = writeTo(data, hmac)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0644)
}

// Resource enums.
const (
	KindCron      = "cron"
	KindPipeline  = "pipeline"
	KindRegistry  = "registry"
	KindSecret    = "secret"
	KindSignature = "signature"
)

type (
	// Manifest is a collection of Drone resources.
	Manifest struct {
		Resources []Resource
	}

	// Resource represents a Drone resource.
	Resource interface {
		// GetVersion returns the resource version.
		GetVersion() string

		// GetKind returns the resource kind.
		GetKind() string
	}

	// RawResource is a raw encoded resource with the
	// resource kind and type extracted.
	RawResource struct {
		Version string
		Kind    string
		Type    string
		Data    []byte `yaml:"-"`
	}

	resource struct {
		Version string
		Kind    string `json:"kind"`
		Type    string `json:"type"`
	}
)

func writeTo(data []byte, hmac string) ([]byte, error) {
	res, err := parseRawBytes(data)
	return upsert(res, hmac), err
}

func parseRawBytes(b []byte) ([]*RawResource, error) {
	return parseRaw(
		bytes.NewReader(b),
	)
}

func parseRaw(r io.Reader) ([]*RawResource, error) {
	const newline = '\n'
	var resources []*RawResource
	var resource *RawResource

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if isSeparator(line) {
			resource = nil
		}
		if resource == nil {
			resource = &RawResource{}
			resources = append(resources, resource)
		}
		if isSeparator(line) {
			continue
		}
		if isTerminator(line) {
			break
		}
		if scanner.Err() == io.EOF {
			break
		}
		resource.Data = append(
			resource.Data,
			line...,
		)
		resource.Data = append(
			resource.Data,
			newline,
		)
	}
	for _, resource := range resources {
		err := yaml.Unmarshal(resource.Data, resource)
		if err != nil {
			return nil, err
		}
	}
	return resources, nil
}

func upsert(res []*RawResource, hmac string) []byte {
	var buf bytes.Buffer
	for _, r := range res {
		if r.Kind != KindSignature {
			buf.WriteString("---")
			buf.WriteByte('\n')
			buf.Write(r.Data)
		}
	}
	buf.WriteString("---")
	buf.WriteByte('\n')
	buf.WriteString("kind: signature")
	buf.WriteByte('\n')
	buf.WriteString("hmac: " + hmac)
	buf.WriteByte('\n')
	buf.WriteByte('\n')
	buf.WriteString("...")
	buf.WriteByte('\n')
	return buf.Bytes()
}

func isSeparator(s string) bool {
	return strings.HasPrefix(s, "---")
}

func isTerminator(s string) bool {
	return strings.HasPrefix(s, "...")
}
