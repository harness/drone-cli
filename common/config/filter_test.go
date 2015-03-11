package config

import (
	"testing"

	"github.com/drone/drone-cli/common"
	"github.com/franela/goblin"
)

var yamlfilter = `
publish:
  google:
    image: google_storage
  docker:
    image: docker
    when:
      branch: dev

deploy:
  amazon:
    image: aws

  microsoft:
    image: azure
    when:
      branch: dev

notify:
  email:
    image: email
    when:
      branch: dev
  slack:
    image: slack

`

func Test_Filter(t *testing.T) {

	repo := &common.Repo{}
	commit := &common.Commit{Branch: "master"}
	conf, err := Parse(yamlfilter)
	if err != nil {
		t.Fail()
	}
	Filter(conf, repo, commit)

	g := goblin.Goblin(t)
	g.Describe("Filter", func() {

		g.It("Should remove steps that don't match condition", func() {
			g.Assert(len(conf.Publish)).Equal(1)
			g.Assert(len(conf.Deploy)).Equal(1)
			g.Assert(len(conf.Notify)).Equal(1)
			g.Assert(conf.Publish["google"].Image).Equal("google_storage")
			g.Assert(conf.Deploy["amazon"].Image).Equal("aws")
			g.Assert(conf.Notify["slack"].Image).Equal("slack")
		})
	})
}
