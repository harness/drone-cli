// +build ignore

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	var files []string
	err := filepath.Walk("files", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".libsonnet" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	buf.WriteString("package stdlib\n\n")
	buf.WriteString("import jsonnet \"github.com/google/go-jsonnet\"\n\n")
	buf.WriteString("var files = map[string]jsonnet.Contents{\n")
	for _, file := range files {
		raw, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(buf, "\t%q:jsonnet.MakeContents(%q),\n", strings.TrimPrefix(file, "files/"), string(raw))
	}
	buf.WriteString("}\n")

	formatted, err := format(buf)
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(formatted)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile("stdlib_gen.go", data, 0644)
}

// format formats a template using gofmt.
func format(in io.Reader) (io.Reader, error) {
	var out bytes.Buffer

	gofmt := exec.Command("gofmt", "-s")
	gofmt.Stdin = in
	gofmt.Stdout = &out
	gofmt.Stderr = os.Stderr
	err := gofmt.Run()
	return &out, err
}
