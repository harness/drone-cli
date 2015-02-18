package config

import (
	"testing"
)

func Test_Parse(t *testing.T) {

	var raw = `
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

	confs, err := ParseMatrix(raw)
	if err != nil {
		t.Error(err)
		return
	}

	if len(confs) != 24 {
		t.Errorf("Expected 24 permutations in matrix, got %d", len(confs))
	}

	unique := map[string]bool{}
	for _, config := range confs {
		unique[config.Image] = true
	}

	if len(unique) != 24 {
		t.Errorf("Expected 24 unique permutations in matrix, got %d", len(unique))
	}
}

/*
	//var axis = map[string][]string{}

	// calculate number of permutations and
	// extract the list of keys.
	var perm int
	var keys []string
	for k, v := range m {
		perm *= len(v)
		if perm == 0 {
			perm = len(v)
		}
		keys = append(keys, k)
	}

	//axisList := map[int][]string{}
	for p := 0; p < perm; p++ {
		var axis []string
		var pos = perm
		for _, key := range keys {
			vals := m[key]
			pos = pos / len(vals)
			index := p / pos % len(vals)
			axis = append(axis, vals[index])
		}
		//axisList[p] = axis
		fmt.Println(axis)
	}

	//fmt.Println(perm)
	//for p := 0; p < perm; p++ {
	//	axis := axisList[p]
	//	data, _ := json.Marshal(axis)
	//	fmt.Println(p, string(data))
	//}

}
*/
