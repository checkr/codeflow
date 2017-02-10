package kubedeploy_test

import (
	"testing"

	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
	"github.com/checkr/codeflow/server/plugins/kubedeploy/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestDeployments struct {
	suite.Suite
	agent *agent.Agent
}

var testDeploymentsConfig = []byte(`
plugins:
  kubedeploy:
    workers: 1
`)

func (suite *TestDeployments) SetupSuite() {
	ag, _ := agent.NewTestAgent(testDeploymentsConfig)
	suite.agent = ag
	go suite.agent.Run()
}

func (suite *TestDeployments) TearDownSuite() {
	suite.agent.Stop()
}

func (suite *TestDeployments) TestSuccessfulDeployment() {
	var e agent.Event

	suite.agent.Events <- testdata.CreateSuccessDeploy()
	e = suite.agent.GetTestEvent("plugins.DockerDeploy:status", 120)
	assert.Equal(suite.T(), string(plugins.Running), string(e.Payload.(plugins.DockerDeploy).State))

	e = suite.agent.GetTestEvent("plugins.DockerDeploy:status", 120)
	assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.DockerDeploy).State))
	for _, service := range e.Payload.(plugins.DockerDeploy).Services {
		assert.Equal(suite.T(), string(plugins.Complete), string(service.State))
	}
}

func (suite *TestDeployments) TestFailedDeployment() {
	var e agent.Event

	suite.agent.Events <- testdata.CreateFailDeploy()
	e = suite.agent.GetTestEvent("plugins.DockerDeploy:status", 120)
	assert.Equal(suite.T(), string(plugins.Running), string(e.Payload.(plugins.DockerDeploy).State))

	e = suite.agent.GetTestEvent("plugins.DockerDeploy:status", 120)
	assert.Equal(suite.T(), string(plugins.Failed), string(e.Payload.(plugins.DockerDeploy).State))
	for _, service := range e.Payload.(plugins.DockerDeploy).Services {
		assert.Equal(suite.T(), string(plugins.Failed), string(service.State))
	}
}

func (suite *TestDeployments) TestStragglerDeployment() {
	var e agent.Event
	// Create a successful deploy
	suite.agent.Events <- testdata.CreateSuccessDeploy()
	// Consume the in-progress message
	e = suite.agent.GetTestEvent("plugins.DockerDeploy:status", 120)
	assert.Equal(suite.T(), string(plugins.Running), string(e.Payload.(plugins.DockerDeploy).State))
	// Consume the success message
	e = suite.agent.GetTestEvent("plugins.DockerDeploy:status", 120)
	assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.DockerDeploy).State))
	// Rename the services and re-deploy
	suite.agent.Events <- testdata.CreateSuccessDeployRenamed()
	// Consume the in-progress message
	e = suite.agent.GetTestEvent("plugins.DockerDeploy:status", 120)
	assert.Equal(suite.T(), string(plugins.Running), string(e.Payload.(plugins.DockerDeploy).State))
	// Consume the success message
	e = suite.agent.GetTestEvent("plugins.DockerDeploy:status", 120)
	assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.DockerDeploy).State))
}

func TestKubeDeployDeployments(t *testing.T) {
	suite.Run(t, new(TestDeployments))
}
