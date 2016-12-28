package kubedeploy_test

import (
	"testing"
	"time"

	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
	"github.com/checkr/codeflow/server/plugins/kubedeploy/testdata"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
}

func (suite *TestSuite) TestKubeDeployDeployments() {
	ag := agent.Agent{}
	shutdown := make(chan struct{})

	// channel shared between all plugin threads for accumulating events
	ag.Events = make(chan agent.Event, 10000)

	ag.TestSubscribeTo = []string{
		"plugins.DockerDeploy:status",
	}

	ag.TestWork = func(e *agent.Event, eventCount int) {
		spew.Dump(eventCount, e.Payload.(plugins.DockerDeploy).State)
		switch eventCount {
		case 1:
			// Deploy in-progress
			assert.Equal(suite.T(), string(plugins.Running), string(e.Payload.(plugins.DockerDeploy).State))
		case 2:
			// Test successful deploy is successful
			assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.DockerDeploy).State))
			// Check that individual service was also marked successful
			for _, service := range e.Payload.(plugins.DockerDeploy).Services {
				assert.Equal(suite.T(), string(plugins.Complete), string(service.State))
			}
		case 3:
			// Failure Deploy in-progress
			assert.Equal(suite.T(), string(plugins.Running), string(e.Payload.(plugins.DockerDeploy).State))
		case 4:
			// Test failure deploy is failure
			assert.Equal(suite.T(), string(plugins.Failed), string(e.Payload.(plugins.DockerDeploy).State))
			for _, service := range e.Payload.(plugins.DockerDeploy).Services {
				assert.Equal(suite.T(), string(plugins.Failed), string(service.State))
			}
			close(shutdown)
			/*
				case 5:
					// In-progress for mixed deploy
					assert.Equal(suite.T(), string(plugins.Running), string(e.Payload.(plugins.DockerDeploy).State))
				case 6:
					// Success mixed deploy
					assert.Equal(suite.T(), string(plugins.Failed), string(e.Payload.(plugins.DockerDeploy).State))
					for _, service := range e.Payload.(plugins.DockerDeploy).Services {
						assert.Equal(suite.T(), string(plugins.Complete), string(service.State))
					}
				case 7:
					// In-progress on teardown
				case 8:
					// Test teardown complete
					close(shutdown)
			*/
		}
	}

	// adding plugin kubedeploy
	creator, _ := agent.PluginRegistry["kubedeploy"]
	plugin := creator()
	rp := &agent.RunningPlugin{
		Name:    "kubedeploy",
		Plugin:  plugin,
		Enabled: true,
	}
	ag.Plugins = append(ag.Plugins, rp)

	// Sequence of test events that trigger all the actions required.
	//testdata.TearDownPreviousDeploys(ag)
	ag.Events <- testdata.CreateSuccessDeploy()
	ag.Events <- testdata.CreateFailDeploy()
	//ag.Events <- testdata.CreateSuccessMixedActionDeploy()

	// Cleanup
	//testdata.TearDownPreviousDeploys(ag)

	// global timeout in the case that we don't ever get all our messages
	timer := time.NewTimer(time.Second * 300)

	go func() {
		<-timer.C
		println("Timer expired, sending exit.")
		assert.Equal(suite.T(), true, timer.Stop())
		close(shutdown)
	}()

	ag.RunTest(shutdown)
}

func TestKubeDeployDeployments(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
