package config

import (
	"github.com/franela/goblin"
	"testing"
)

func Test_Inject(t *testing.T) {

	g := goblin.Goblin(t)
	g.Describe("Inject params", func() {

		g.It("Should replace vars with $$", func() {
			s := "echo $$FOO $BAR"
			m := map[string]string{}
			m["FOO"] = "BAZ"
			g.Assert("echo BAZ $BAR").Equal(Inject(s, m))
		})

		g.It("Should not replace vars with single $", func() {
			s := "echo $FOO $BAR"
			m := map[string]string{}
			m["FOO"] = "BAZ"
			g.Assert(s).Equal(Inject(s, m))
		})

		g.It("Should not replace vars in nil map", func() {
			s := "echo $$FOO $BAR"
			g.Assert(s).Equal(Inject(s, nil))
		})
	})
}

func Test_InjectSafe(t *testing.T) {

	g := goblin.Goblin(t)
	g.Describe("Safely Inject params", func() {

		m := map[string]string{}
		m["TOKEN"] = "FOO"
		m["SECRET"] = "BAR"
		c, _ := Parse(InjectSafe(yml, m))

		g.It("Should replace vars in notify section", func() {
			g.Assert(c.Deploy["my_service"].(map[interface{}]interface{})["token"]).Equal("FOO")
			g.Assert(c.Deploy["my_service"].(map[interface{}]interface{})["secret"]).Equal("BAR")
		})

		g.It("Should not replace vars in script section", func() {
			g.Assert(c.Script[0]).Equal("echo $$TOKEN")
			g.Assert(c.Script[1]).Equal("echo $$SECRET")
		})
	})
}

var yml = `
image: foo
script:
  - echo $$TOKEN
  - echo $$SECRET
deploy:
  my_service:
    token: $$TOKEN
    secret: $$SECRET
`
