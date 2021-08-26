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
	}{
		{
			name:        "Stream + Format",
			jsonnetFile: "stream_format.jsonnet",
			yamlFile:    "stream_format.yaml",
			format:      true, stream: true,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			expected, err := os.ReadFile(filepath.Join("./testdata", tc.yamlFile))
			assert.NoError(t, err)

			result, err := convert(filepath.Join("./testdata", tc.jsonnetFile), tc.stringOutput, tc.format, tc.stream, tc.extVars)
			assert.NoError(t, err)
			assert.Equal(t, string(expected), result)
		})
	}
}
