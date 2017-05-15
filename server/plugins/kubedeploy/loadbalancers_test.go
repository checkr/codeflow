package kubedeploy_test

import (
	"testing"

	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
	"github.com/checkr/codeflow/server/plugins/kubedeploy/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestLoadBalancers struct {
	suite.Suite
	agent *agent.Agent
}

var lbConfig = []byte(`
plugins:
  kubedeploy:
    workers: 1
`)

func (suite *TestLoadBalancers) SetupSuite() {
	ag, _ := agent.NewTestAgent(testDeploymentsConfig)
	suite.agent = ag
	go suite.agent.Run()

	// Test teardown of any existing LBs
	suite.agent.Events <- testdata.TearDownLBTCP(plugins.Internal)
	_ = suite.agent.GetTestEvent("plugins.LoadBalancer:status", 60)
	suite.agent.Events <- testdata.TearDownLBHTTPS(plugins.Internal)
	_ = suite.agent.GetTestEvent("plugins.LoadBalancer:status", 60)
	suite.agent.Events <- testdata.TearDownLBTCP(plugins.External)
	_ = suite.agent.GetTestEvent("plugins.LoadBalancer:status", 60)
	suite.agent.Events <- testdata.TearDownLBHTTPS(plugins.External)
	_ = suite.agent.GetTestEvent("plugins.LoadBalancer:status", 60)
}

func (suite *TestLoadBalancers) TearDownSuite() {
	suite.agent.Stop()
}

func (suite *TestLoadBalancers) TestLBTCPInternal() {
	var e agent.Event

	suite.agent.Events <- testdata.CreateLBTCP(plugins.Internal)
	e = suite.agent.GetTestEvent("plugins.LoadBalancer:status", 60)
	assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.LoadBalancer).State))

	suite.agent.Events <- testdata.UpdateLBTCP(plugins.Internal)
	e = suite.agent.GetTestEvent("plugins.LoadBalancer:status", 60)
	assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.LoadBalancer).State))

	suite.agent.Events <- testdata.TearDownLBTCP(plugins.Internal)
	e = suite.agent.GetTestEvent("plugins.LoadBalancer:status", 60)
	assert.Equal(suite.T(), string(plugins.Deleted), string(e.Payload.(plugins.LoadBalancer).State))
}

func (suite *TestLoadBalancers) TestLBHTTPSInternal() {
	var e agent.Event

	suite.agent.Events <- testdata.CreateLBHTTPS(plugins.Internal)
	e = suite.agent.GetTestEvent("plugins.LoadBalancer:status", 60)
	assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.LoadBalancer).State))

	// Test updating LBHTTPS LB
	suite.agent.Events <- testdata.UpdateLBHTTPS(plugins.Internal)
	e = suite.agent.GetTestEvent("plugins.LoadBalancer:status", 60)
	assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.LoadBalancer).State))

	suite.agent.Events <- testdata.TearDownLBHTTPS(plugins.Internal)
	e = suite.agent.GetTestEvent("plugins.LoadBalancer:status", 60)
	assert.Equal(suite.T(), string(plugins.Deleted), string(e.Payload.(plugins.LoadBalancer).State))
}

func (suite *TestLoadBalancers) TestLBTCPExternal() {
	var e agent.Event

	suite.agent.Events <- testdata.CreateLBTCP(plugins.External)
	e = suite.agent.GetTestEvent("plugins.LoadBalancer:status", 60)
	assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.LoadBalancer).State))
	assert.NotNil(suite.T(), string(e.Payload.(plugins.LoadBalancer).DNS))

	suite.agent.Events <- testdata.UpdateLBTCP(plugins.External)
	e = suite.agent.GetTestEvent("plugins.LoadBalancer:status", 60)
	assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.LoadBalancer).State))
	assert.NotNil(suite.T(), string(e.Payload.(plugins.LoadBalancer).DNS))

	suite.agent.Events <- testdata.TearDownLBTCP(plugins.External)
	e = suite.agent.GetTestEvent("plugins.LoadBalancer:status", 60)
	assert.Equal(suite.T(), string(plugins.Deleted), string(e.Payload.(plugins.LoadBalancer).State))
}

func (suite *TestLoadBalancers) TestLBHTTPSExternal() {
	var e agent.Event

	suite.agent.Events <- testdata.CreateLBHTTPS(plugins.External)
	e = suite.agent.GetTestEvent("plugins.LoadBalancer:status", 60)
	assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.LoadBalancer).State))
	assert.NotNil(suite.T(), string(e.Payload.(plugins.LoadBalancer).DNS))

	suite.agent.Events <- testdata.UpdateLBHTTPS(plugins.External)
	e = suite.agent.GetTestEvent("plugins.LoadBalancer:status", 60)
	assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.LoadBalancer).State))
	assert.NotNil(suite.T(), string(e.Payload.(plugins.LoadBalancer).DNS))

	suite.agent.Events <- testdata.TearDownLBHTTPS(plugins.External)
	e = suite.agent.GetTestEvent("plugins.LoadBalancer:status", 60)
	assert.Equal(suite.T(), string(plugins.Deleted), string(e.Payload.(plugins.LoadBalancer).State))
}

func TestKubeDeployLoadbalancers(t *testing.T) {
	suite.Run(t, new(TestLoadBalancers))
}
