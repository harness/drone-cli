package secure

import (
	"testing"

	"github.com/franela/goblin"
	"gopkg.in/yaml.v2"
)

func Test_MapEqualSlice(t *testing.T) {

	g := goblin.Goblin(t)
	g.Describe("MapEqualSlice", func() {

		g.It("Should unmarshal map", func() {

			out := &Secure{}
			err := yaml.Unmarshal([]byte(mapYaml), out)
			g.Assert(err == nil).IsTrue()
			g.Assert(out.Environment.Map()["FOO"]).Equal("BAR")
			g.Assert(out.Environment.Map()["BAZ"]).Equal("BOO")

		})

		g.It("Should unmarshal slice", func() {
			out := &Secure{}
			err := yaml.Unmarshal([]byte(sliceYaml), out)
			g.Assert(err == nil).IsTrue()
			g.Assert(out.Environment.Map()["FOO"]).Equal("BAR")
			g.Assert(out.Environment.Map()["BAZ"]).Equal("BOO")
		})
	})
}
