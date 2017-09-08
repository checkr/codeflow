package kubedeploy_test

import (
	"testing"
	"time"

	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
	"github.com/checkr/codeflow/server/plugins/kubedeploy/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestDeployments struct {
	suite.Suite
	agent1 *agent.Agent
	agent2 *agent.Agent
}

var testDeploymentsConfig = []byte(`
plugins:
  kubedeploy:
    workers: 1
`)

func (suite *TestDeployments) SetupSuite() {
	ag1, _ := agent.NewTestAgent(testDeploymentsConfig)
	ag2, _ := agent.NewTestAgent(testDeploymentsConfig)

	suite.agent1 = ag1
	suite.agent2 = ag2

	go suite.agent1.Run()
	go suite.agent2.Run()
}

func (suite *TestDeployments) TearDownSuite() {
	suite.agent1.Stop()
	suite.agent2.Stop()
}

func (suite *TestDeployments) TestSuccessfulDeployment() {
	var e agent.Event

	suite.agent1.Events <- testdata.CreateSuccessDeploy()
	e = suite.agent1.GetTestEvent("plugins.DockerDeploy:status", 120)
	assert.Equal(suite.T(), string(plugins.Running), string(e.Payload.(plugins.DockerDeploy).State))

	e = suite.agent1.GetTestEvent("plugins.DockerDeploy:status", 120)
	assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.DockerDeploy).State))
	for _, service := range e.Payload.(plugins.DockerDeploy).Services {
		assert.Equal(suite.T(), string(plugins.Complete), string(service.State))
	}
}

func (suite *TestDeployments) TestSuccessfulJob() {
	var e agent.Event

	suite.agent1.Events <- testdata.CreateSuccessJob()

	e = suite.agent1.GetTestEvent("plugins.DockerDeploy:status", 120)
	assert.Equal(suite.T(), string(plugins.Running), string(e.Payload.(plugins.DockerDeploy).State))
	for _, service := range e.Payload.(plugins.DockerDeploy).Services {
		assert.Equal(suite.T(), true, service.OneShot)
		assert.Equal(suite.T(), string(plugins.Waiting), string(service.State))
	}

	e = suite.agent1.GetTestEvent("plugins.DockerDeploy:status", 120)
	assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.DockerDeploy).State))
	for _, service := range e.Payload.(plugins.DockerDeploy).Services {
		assert.Equal(suite.T(), true, service.OneShot)
		assert.Equal(suite.T(), string(plugins.Complete), string(service.State))
	}
}

func (suite *TestDeployments) TestFailedJobImagePull() {
	var e agent.Event

	suite.agent1.Events <- testdata.CreateFailJob()
	e = suite.agent1.GetTestEvent("plugins.DockerDeploy:status", 120)

	for _, service := range e.Payload.(plugins.DockerDeploy).Services {
		assert.Equal(suite.T(), true, service.OneShot)
		assert.Equal(suite.T(), string(plugins.Waiting), string(service.State))
	}

	e = suite.agent1.GetTestEvent("plugins.DockerDeploy:status", 120)
	for _, service := range e.Payload.(plugins.DockerDeploy).Services {
		assert.Equal(suite.T(), true, service.OneShot)
		assert.Equal(suite.T(), string(plugins.Failed), string(service.State))
	}
}

func (suite *TestDeployments) TestFailedJob() {
	var e agent.Event

	suite.agent1.Events <- testdata.CreateFailJobNonZero()

	e = suite.agent1.GetTestEvent("plugins.DockerDeploy:status", 120)

	for _, service := range e.Payload.(plugins.DockerDeploy).Services {
		assert.Equal(suite.T(), true, service.OneShot)
		assert.Equal(suite.T(), string(plugins.Waiting), string(service.State))

	}

	e = suite.agent1.GetTestEvent("plugins.DockerDeploy:status", 120)
	for _, service := range e.Payload.(plugins.DockerDeploy).Services {
		assert.Equal(suite.T(), true, service.OneShot)
		assert.Equal(suite.T(), string(plugins.Failed), string(service.State))
	}
}

func (suite *TestDeployments) TestFailIfJobAlreadyActive() {
	var e agent.Event

	suite.agent1.Events <- testdata.CreateAlreadyActiveSoFailJob("SERVICE-JOB-1")
	time.Sleep(time.Second * 5)
	suite.agent2.Events <- testdata.CreateAlreadyActiveSoFailJob("SERVICE-JOB-1")

	e = suite.agent1.GetTestEvent("plugins.DockerDeploy:status", 120)
	for _, service := range e.Payload.(plugins.DockerDeploy).Services {
		assert.Equal(suite.T(), string(plugins.Waiting), string(service.State))
	}

	e = suite.agent2.GetTestEvent("plugins.DockerDeploy:status", 120)
	for _, service := range e.Payload.(plugins.DockerDeploy).Services {
		assert.Equal(suite.T(), string(plugins.Waiting), string(service.State))
	}

	e = suite.agent2.GetTestEvent("plugins.DockerDeploy:status", 120)
	for _, service := range e.Payload.(plugins.DockerDeploy).Services {
		assert.Equal(suite.T(), string(plugins.Failed), string(service.State))
	}

	e = suite.agent1.GetTestEvent("plugins.DockerDeploy:status", 120)
	for _, service := range e.Payload.(plugins.DockerDeploy).Services {
		assert.Equal(suite.T(), string(plugins.Complete), string(service.State))
	}

}

func (suite *TestDeployments) TestFailedDeploymentImagePull() {
	var e agent.Event

	suite.agent1.Events <- testdata.CreateFailDeploy()
	e = suite.agent1.GetTestEvent("plugins.DockerDeploy:status", 120)
	assert.Equal(suite.T(), string(plugins.Running), string(e.Payload.(plugins.DockerDeploy).State))

	e = suite.agent1.GetTestEvent("plugins.DockerDeploy:status", 120)
	assert.Equal(suite.T(), string(plugins.Failed), string(e.Payload.(plugins.DockerDeploy).State))
	for _, service := range e.Payload.(plugins.DockerDeploy).Services {
		assert.Equal(suite.T(), string(plugins.Failed), string(service.State))
	}
}

func (suite *TestDeployments) TestFailedDeploymentCommand() {
	var e agent.Event

	suite.agent1.Events <- testdata.CreateFailDeployCommand()
	e = suite.agent1.GetTestEvent("plugins.DockerDeploy:status", 120)
	assert.Equal(suite.T(), string(plugins.Running), string(e.Payload.(plugins.DockerDeploy).State))

	e = suite.agent1.GetTestEvent("plugins.DockerDeploy:status", 120)
	assert.Equal(suite.T(), string(plugins.Failed), string(e.Payload.(plugins.DockerDeploy).State))
	for _, service := range e.Payload.(plugins.DockerDeploy).Services {
		// Check if service is one shot
		if service.OneShot == true {
			assert.Equal(suite.T(), string(plugins.Failed), string(service.State))
		}
		assert.Equal(suite.T(), string(plugins.Failed), string(service.State))
	}
}

func (suite *TestDeployments) TestStragglerDeployment() {
	var e agent.Event
	// Create a successful deploy
	suite.agent1.Events <- testdata.CreateSuccessDeploy()
	// Consume the in-progress message
	e = suite.agent1.GetTestEvent("plugins.DockerDeploy:status", 120)
	assert.Equal(suite.T(), string(plugins.Running), string(e.Payload.(plugins.DockerDeploy).State))
	// Consume the success message
	e = suite.agent1.GetTestEvent("plugins.DockerDeploy:status", 120)
	assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.DockerDeploy).State))
	// Rename the services and re-deploy
	suite.agent1.Events <- testdata.CreateSuccessDeployRenamed()
	// Consume the in-progress message
	e = suite.agent1.GetTestEvent("plugins.DockerDeploy:status", 120)
	assert.Equal(suite.T(), string(plugins.Running), string(e.Payload.(plugins.DockerDeploy).State))
	// Consume the success message
	e = suite.agent1.GetTestEvent("plugins.DockerDeploy:status", 120)
	assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.DockerDeploy).State))
}

func (suite *TestDeployments) TestFailedDeployFollowedBySuccessDeploy() {
	// To exercise the fast failure mode we will make sure that the detector picks the right replica set.
	var e agent.Event

	// Create a failed deploy with pods waiting 'forever' condition.
	suite.agent1.Events <- testdata.CreateSuccessAndFailDeploy1()
	// Consume the in-progress message
	e = suite.agent1.GetTestEvent("plugins.DockerDeploy:status", 120)
	assert.Equal(suite.T(), string(plugins.Running), string(e.Payload.(plugins.DockerDeploy).State))
	// Consume the success message
	e = suite.agent1.GetTestEvent("plugins.DockerDeploy:status", 120)
	assert.Equal(suite.T(), string(plugins.Failed), string(e.Payload.(plugins.DockerDeploy).State))

	// Create a success deploy and make sure it succeeds.
	suite.agent1.Events <- testdata.CreateSuccessAndFailDeploy2()
	// Consume the in-progress message
	e = suite.agent1.GetTestEvent("plugins.DockerDeploy:status", 120)
	assert.Equal(suite.T(), string(plugins.Running), string(e.Payload.(plugins.DockerDeploy).State))
	// Consume the success message
	e = suite.agent1.GetTestEvent("plugins.DockerDeploy:status", 120)
	assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.DockerDeploy).State))
}

func TestKubeDeployDeployments(t *testing.T) {
	suite.Run(t, new(TestDeployments))
}
