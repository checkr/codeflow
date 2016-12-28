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

type TestSuiteLB struct {
	suite.Suite
}

func (suite *TestSuiteLB) TestKubeDeployLoadbalancers() {
	ag := agent.Agent{}
	shutdown := make(chan struct{})

	// channel shared between all plugin threads for accumulating events
	ag.Events = make(chan agent.Event, 10000)

	ag.TestSubscribeTo = []string{
		"plugins.LoadBalancer:status",
	}

	ag.TestWork = func(e *agent.Event, eventCount int) {
		switch eventCount {
		case 1:
			// Test teardown of any existing LBs
			// discard first response
		case 2:
			// Test teardown of any existing LBs
			// discard first response
		case 3:
			// Test creating LBTCP
			assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.LoadBalancer).State))
		case 4:
			// Test updating LBTCP
			assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.LoadBalancer).State))
		case 5:
			// Test creating LBHTTPS LB
			assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.LoadBalancer).State))
		case 6:
			// Test updating LBHTTPS LB
			assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.LoadBalancer).State))
		case 7:
			// Test Teardown Success / Destroy
			assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.LoadBalancer).State))
		case 8:
			// Test Teardown Success / Destroy
			assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.LoadBalancer).State))
		case 9:
			// Test teardown of any existing LBs
			// discard first response
		case 10:
			// Test teardown of any existing LBs
			// discard first response
		case 11:
			// Test creating LBTCP
			assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.LoadBalancer).State))
			assert.NotNil(suite.T(), string(e.Payload.(plugins.LoadBalancer).DNSName))
		case 12:
			// Test updating LBTCP LB
			assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.LoadBalancer).State))
			assert.NotNil(suite.T(), string(e.Payload.(plugins.LoadBalancer).DNSName))
		case 13:
			// Test creating LBHTTPS
			assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.LoadBalancer).State))
			assert.NotNil(suite.T(), string(e.Payload.(plugins.LoadBalancer).DNSName))
		case 14:
			// Test updating LBHTTPS LB
			assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.LoadBalancer).State))
			assert.NotNil(suite.T(), string(e.Payload.(plugins.LoadBalancer).DNSName))
		case 15:
			// Test Teardown Success / Destroy
			assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.LoadBalancer).State))
		case 16:
			// Test Teardown Success / Destroy
			assert.Equal(suite.T(), string(plugins.Complete), string(e.Payload.(plugins.LoadBalancer).State))
			close(shutdown)
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
	ag.Events <- testdata.TearDownLBTCP(plugins.Internal)
	ag.Events <- testdata.TearDownLBHTTPS(plugins.Internal)
	ag.Events <- testdata.CreateLBTCP(plugins.Internal)
	ag.Events <- testdata.UpdateLBTCP(plugins.Internal)
	ag.Events <- testdata.CreateLBHTTPS(plugins.Internal)
	ag.Events <- testdata.UpdateLBHTTPS(plugins.Internal)
	// use teardown to test LB Destroy action
	ag.Events <- testdata.TearDownLBTCP(plugins.Internal)
	ag.Events <- testdata.TearDownLBHTTPS(plugins.Internal)

	ag.Events <- testdata.TearDownLBTCP(plugins.External)
	ag.Events <- testdata.TearDownLBHTTPS(plugins.External)
	ag.Events <- testdata.CreateLBTCP(plugins.External)
	ag.Events <- testdata.UpdateLBTCP(plugins.External)
	ag.Events <- testdata.CreateLBHTTPS(plugins.External)
	ag.Events <- testdata.UpdateLBHTTPS(plugins.External)
	// use teardown to test LB Destroy action
	ag.Events <- testdata.TearDownLBTCP(plugins.External)
	ag.Events <- testdata.TearDownLBHTTPS(plugins.External)

	// global timeout in the case that we don't ever get all our messages
	timer := time.NewTimer(time.Second * 120)

	go func() {
		<-timer.C
		println("Timer expired, sending exit.")
		assert.Equal(suite.T(), true, timer.Stop())
		close(shutdown)
	}()

	ag.RunTest(shutdown)
}

func TestKubeDeployLoadbalancers(t *testing.T) {
	suite.Run(t, new(TestSuiteLB))
}
