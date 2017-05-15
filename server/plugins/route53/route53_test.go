package route53_test

import (
	"log"
	"strings"
	"testing"

	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
	"github.com/checkr/codeflow/server/plugins/kubedeploy/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	_ "github.com/checkr/codeflow/server/plugins/kubedeploy"
	"github.com/spf13/viper"
)

type TestRoute53 struct {
	suite.Suite
	agent *agent.Agent
}

var testRoute53Config = []byte(`
plugins:
  kubedeploy:
    workers: 1
  route53:
    workers: 1
`)

func (suite *TestRoute53) SetupSuite() {
	ag, _ := agent.NewTestAgent(testRoute53Config)
	suite.agent = ag
	go suite.agent.Run()
}

func (suite *TestRoute53) TearDownSuite() {
	suite.agent.Stop()
}

func TestRoute53Plugin(t *testing.T) {
	suite.Run(t, new(TestRoute53))
}

func (suite *TestRoute53) TestRoute53Update() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("CF")
	viper.AutomaticEnv()

	if !viper.IsSet("plugins.route53.hosted_zone_id") || !viper.IsSet("plugins.route53.hosted_zone_name") {
		log.Println("You must set CF_PLUGINS_ROUTE53_HOSTED_ZONE_ID and CF_PLUGINS_ROUTE53_HOSTED_ZONE_NAME to run this test.")
	}
	var e agent.Event
	suite.agent.Events <- testdata.CreateLBHTTPS(plugins.External)

	// Route53 message
	e = suite.agent.GetTestEvent("plugins.Route53", 900)
	route53Response := e.Payload.(plugins.Route53)
	assert.Equal(suite.T(), string(plugins.Complete), string(route53Response.State))
}
