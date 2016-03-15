package drone

import (
	"encoding/json"
	"testing"
)

func TestNetrcMarshaling(t *testing.T) {
	var tests = []struct {
		num    int
		input  string
		output string
	}{
		{
			1,
			`{"machine": "Apple 2", "login": "MeMyselfAndI", "password": "pass123"}`,
			`{"machine":"Apple 2","login":"MeMyselfAndI","password":"pass123","user":"pass123"}`,
		},
		{
			2,
			`{"machine": "Apple 2", "login": "MeMyselfAndI", "user": "pass123"}`,
			`{"machine":"Apple 2","login":"MeMyselfAndI","password":"pass123","user":"pass123"}`,
		},
		{
			3,
			`{"machine": "Apple 2", "login": "MeMyselfAndI"}`,
			`{"machine":"Apple 2","login":"MeMyselfAndI","password":"","user":""}`,
		},
		{
			4,
			`{"machine": "Apple 2", "password": "pass123"}`,
			`{"machine":"Apple 2","login":"","password":"pass123","user":"pass123"}`,
		},
		{
			5,
			`{"machine": "Apple 2", "user": "pass123"}`,
			`{"machine":"Apple 2","login":"","password":"pass123","user":"pass123"}`,
		},
		{
			6,
			`{ "login": "MeMyselfAndI", "password": "pass123"}`,
			`{"machine":"","login":"MeMyselfAndI","password":"pass123","user":"pass123"}`,
		},
		{
			7,
			`{ "login": "MeMyselfAndI", "user": "pass123"}`,
			`{"machine":"","login":"MeMyselfAndI","password":"pass123","user":"pass123"}`,
		},
	}

	for _, test := range tests {
		var netrc = Netrc{}
		if err := json.Unmarshal([]byte(test.input), &netrc); err != nil {
			panic(err)
		}
		bytes, err := json.Marshal(netrc)
		if err != nil {
			panic(err)
		}
		if string(bytes) != test.output {
			t.Errorf("test %d expected %q got %q", test.num, test.output, string(bytes))
		}
	}
}
