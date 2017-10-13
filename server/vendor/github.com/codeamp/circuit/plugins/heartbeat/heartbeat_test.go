package heartbeat_test

import (
	"testing"

	"github.com/codeamp/circuit/plugins"
	"github.com/codeamp/transistor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	transistor *transistor.Agent
}

var config = []byte(`
plugins:
  heartbeat:
    workers: 0
`)

func (suite *TestSuite) SetupSuite() {
	ag, _ := transistor.NewTesttransistor(config)
	suite.transistor = ag
	go suite.transistor.Run()
}

func (suite *TestSuite) TearDownSuite() {
	suite.transistor.Stop()
}

func (suite *TestSuite) TestHeartbeat() {
	var e transistor.Event

	e = suite.transistor.GetTestEvent("plugins.HeartBeat", 61)
	payload := e.Payload.(plugins.HeartBeat)
	assert.Equal(suite.T(), "minute", payload.Tick)
}

func TestHeartbeat(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
