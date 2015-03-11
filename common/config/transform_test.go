package config

import (
	"testing"

	"github.com/drone/drone-cli/common"
	"github.com/franela/goblin"
)

func Test_Transform(t *testing.T) {

	g := goblin.Goblin(t)
	g.Describe("Transform", func() {

		g.It("Should add setup step", func() {
			c := &common.Config{}
			c.Build = &common.Step{}
			c.Build.Config = map[string]interface{}{}
			addSetup(c)
			g.Assert(c.Setup != nil).IsTrue()
			g.Assert(c.Setup.Image).Equal("plugins/drone-build")
			g.Assert(c.Setup.Config).Equal(c.Build.Config)
		})

		g.It("Should add clone step", func() {
			c := &common.Config{}
			addClone(c)
			g.Assert(c.Clone != nil).IsTrue()
			g.Assert(c.Clone.Image).Equal("plugins/drone-git")
		})

		g.It("Should normalize build", func() {
			c := &common.Config{}
			c.Build = &common.Step{}
			c.Build.Config = map[string]interface{}{}
			c.Build.Config["commands"] = []string{"echo hello"}
			normalizeBuild(c)
			g.Assert(len(c.Build.Config)).Equal(0)
		})

		g.It("Should normalize images")
		g.It("Should normalize docker plugin")

		g.It("Should remove publish", func() {
			c := &common.Config{}
			c.Publish = map[string]*common.Step{}
			c.Publish["docer"] = &common.Step{}
			rmPublish(c)
			g.Assert(len(c.Publish)).Equal(0)
		})

		g.It("Should remove deploy", func() {
			c := &common.Config{}
			c.Deploy = map[string]*common.Step{}
			c.Deploy["rackspace"] = &common.Step{}
			rmDeploy(c)
			g.Assert(len(c.Deploy)).Equal(0)
		})

		g.It("Should remove notify", func() {
			c := &common.Config{}
			c.Notify = map[string]*common.Step{}
			c.Notify["gmail"] = &common.Step{}
			rmNotify(c)
			g.Assert(len(c.Notify)).Equal(0)
		})

		g.It("Should remove privileged")
		g.It("Should remove volumes")
		g.It("Should remove network")

		g.It("Should return full qualified image name", func() {
			g.Assert(imageName("microsoft/azure")).Equal("microsoft/azure")
			g.Assert(imageName("azure")).Equal("plugins/drone-azure")
			g.Assert(imageName("azure_storage")).Equal("plugins/drone-azure-storage")
		})
	})
}
