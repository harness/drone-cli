package yaml

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Parse parses the configuration from io.Reader r.
func Parse(r io.Reader) (*Manifest, error) {
	manifest := new(Manifest)
	scanner := bufio.NewScanner(r)
	row := 0
	buf := new(bytes.Buffer)
	for scanner.Scan() {
		row++
		txt := scanner.Text()
		if strings.HasPrefix(txt, "---") && row != 1 {
			resource, err := parse(buf.Bytes())
			if err != nil {
				return nil, err
			}
			manifest.Resources = append(
				manifest.Resources,
				resource,
			)
			buf.Reset()
		} else {
			buf.WriteString(txt)
			buf.WriteByte('\n')
		}
	}
	resource, err := parse(buf.Bytes())
	if err != nil {
		return nil, err
	}
	manifest.Resources = append(
		manifest.Resources,
		resource,
	)
	return manifest, nil
}

// ParseBytes parses the configuration from bytes b.
func ParseBytes(b []byte) (*Manifest, error) {
	return Parse(
		bytes.NewBuffer(b),
	)
}

// ParseString parses the configuration from string s.
func ParseString(s string) (*Manifest, error) {
	return ParseBytes(
		[]byte(s),
	)
}

// ParseFile parses the configuration from path p.
func ParseFile(p string) (*Manifest, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}

func parse(b []byte) (Resource, error) {
	res := new(resource)
	err := yaml.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}
	var obj Resource
	switch res.Kind {
	case "cron":
		obj = new(Cron)
	case "secret":
		obj = new(Secret)
	case "signature":
		obj = new(Signature)
	case "registry":
		obj = new(Registry)
	default:
		obj = new(Pipeline)
	}
	err = yaml.Unmarshal(b, obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
