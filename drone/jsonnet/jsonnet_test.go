package jsonnet

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	testcases := []struct {
		name                         string
		jsonnetFile, yamlFile        string
		stringOutput, format, stream bool
		extVars                      []string
		jpath                        []string
	}{
		{
			name:        "Stream + Format",
			jsonnetFile: "stream_format.jsonnet",
			yamlFile:    "stream_format.yaml",
			format:      true, stream: true,
		},
		{
			name:        "Jsonnet Path",
			jsonnetFile: "stream_format.jsonnet",
			yamlFile:    "stream_format.yaml",
			format:      true, stream: true,
			jpath: []string{"/path/to/jsonnet/lib"},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			expected, err := os.ReadFile(filepath.Join("./testdata", tc.yamlFile))
			assert.NoError(t, err)

			result, err := convert(filepath.Join("./testdata", tc.jsonnetFile), tc.stringOutput, tc.format, tc.stream, tc.extVars, tc.jpath)
			assert.NoError(t, err)
			assert.Equal(t, string(expected), result)
		})
	}
}
