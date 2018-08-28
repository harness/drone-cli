package stdlib

import (
	"path"

	jsonnet "github.com/google/go-jsonnet"
)

//go:generate go run gen.go

// Importer provides a default importer that automatically
// loads the embedded drone standard library.
func Importer() jsonnet.Importer {
	return &importer{
		base: &jsonnet.FileImporter{},
	}
}

type importer struct {
	base jsonnet.Importer
}

func (i *importer) Import(importedFrom, importedPath string) (contents jsonnet.Contents, foundAt string, err error) {
	dir, _ := path.Split(importedFrom)
	path := path.Join(dir, importedPath)
	if contents, ok := files[path]; ok {
		return contents, path, nil
	}
	return i.base.Import(importedFrom, importedPath)
}
