package route53_test

import (
	"log"
	"testing"

	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
	"github.com/checkr/codeflow/server/plugins/kubedeploy/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	_ "github.com/checkr/codeflow/server/plugins/kubedeploy"
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
	var e agent.Event
	suite.agent.Events <- testdata.CreateLBHTTPS(plugins.External)
	e = suite.agent.GetTestEvent("plugins.LoadBalancer:status", 60)
	pay := e.Payload.(plugins.LoadBalancer)
	log.Printf("DNSName:%s, State:%s, Message:%s", pay.DNSName, pay.State, pay.StateMessage)
	assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.LoadBalancer).State))
}
