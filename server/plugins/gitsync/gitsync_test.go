package gitsync_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
	log "github.com/codeamp/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var config = []byte(`
plugins:
  gitsync:
    workers: 1
    workdir: "/tmp/gitsync" 
`)

type TestSuite struct {
	suite.Suite
	agent *agent.Agent
}

func (suite *TestSuite) SetupSuite() {
	ag, _ := agent.NewTestAgent(config)
	suite.agent = ag
	go suite.agent.Run()
}

func (suite *TestSuite) TearDownSuite() {
	suite.agent.Stop()
}

func (suite *TestSuite) TestGitsync() {
	var e agent.Event

	log.SetLogLevel("debug")

	gitSync := plugins.GitSync{
		Action: plugins.Update,
		State:  plugins.Waiting,
		Project: plugins.Project{
			Slug:       "codeamp-circuit",
			Repository: "codeamp/circuit",
		},
		Git: plugins.Git{
			Url:           "https://github.com/codeamp/circuit.git",
			Protocol:      "HTTPS",
			Branch:        "master",
			RsaPrivateKey: "",
			RsaPublicKey:  "",
		},
		From: "",
	}

	suite.agent.Events <- agent.NewEvent(gitSync, nil)

	created := time.Now()
	for i := 0; i < 5; i++ {
		e = suite.agent.GetTestEvent("plugins.GitCommit", 60)
		payload := e.Payload.(plugins.GitCommit)
		assert.Equal(suite.T(), payload.Repository, gitSync.Project.Repository)
		assert.Equal(suite.T(), payload.Ref, fmt.Sprintf("refs/heads/%s", gitSync.Git.Branch))
		assert.True(suite.T(), payload.Created.Before(created), "Commit created time is older than previous commit")
	}
}

func TestGitsync(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
