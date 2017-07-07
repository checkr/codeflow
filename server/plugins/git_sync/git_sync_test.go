package git_sync_test

import (
	"testing"

	"github.com/checkr/codeflow/server/agent"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	agent *agent.Agent
}

var config = []byte(`
plugins:
  git_sync:
    workers: 1
`)

func (suite *TestSuite) SetupSuite() {
	ag, _ := agent.NewTestAgent(config)
	suite.agent = ag
	go suite.agent.Run()
}

func (suite *TestSuite) TearDownSuite() {
	suite.agent.Stop()
}

func (suite *TestSuite) TestHeartbeat() {
}

func TestHeartbeat(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
