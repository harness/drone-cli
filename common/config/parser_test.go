package config

import (
	"testing"
)

func Test_Parse(t *testing.T) {
	confs, err := ParseMatrix(matrix)
	if err != nil {
		t.Error(err)
		return
	}

	if len(confs) != 24 {
		t.Errorf("Expected 24 permutations in matrix, got %d", len(confs))
	}

	unique := map[string]bool{}
	for _, config := range confs {
		unique[config.Build.Image] = true
	}

	if len(unique) != 24 {
		t.Errorf("Expected 24 unique permutations in matrix, got %d", len(unique))
	}
}

var matrix = `
build:
  image: $$python_version $$redis_version $$django_version $$go_version
matrix:
  python_version:
    - 3.2
    - 3.3
  redis_version:
    - 2.6
    - 2.8
  django_version:
    - 1.7
    - 1.7.1
    - 1.7.2
  go_version:
    - go1
    - go1.2
`
