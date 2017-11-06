package gitsync_test

import (
	"testing"

	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
	log "github.com/codeamp/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var config = []byte(`
environment: test
plugins:
  dockerbuilder:
    workers: 1
    registry_host: "docker.io"
    registry_org: "checkr"
    registry_username: ""
    registry_password: ""
    registry_user_email: ""
    workdir: "/tmp/dockerbuilder"   
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

func (suite *TestSuite) TestDockerBuilder() {
	var e agent.Event

	log.SetLogLevel(logrus.DebugLevel)

	dockerBuildEvent := plugins.DockerBuild{
		Action: plugins.Create,
		State:  plugins.Waiting,
		Project: plugins.Project{
			Slug:       "checkr-codeflow",
			Repository: "checkr/codeflow",
		},
		Git: plugins.Git{
			Url:           "https://github.com/checkr/codeflow.git",
			Protocol:      "HTTPS",
			Branch:        "master",
			RsaPrivateKey: "",
			RsaPublicKey:  "",
			Workdir:       viper.GetString("plugins.dockerbuilder.workdir"),
		},
		Feature: plugins.Feature{
			Hash:       "51666df1f6d51ec1407796c9645c7c389de0f223",
			ParentHash: "4f059f5d9cd8a65bd4acc6f1d6721ee6eb3fb713",
			User:       "Saso Matejina",
			Message:    "Test",
		},
		Registry: plugins.DockerRegistry{
			Host:     viper.GetString("plugins.dockerbuilder.registry_host"),
			Org:      viper.GetString("plugins.dockerbuilder.registry_org"),
			Username: viper.GetString("plugins.dockerbuilder.registry_username"),
			Password: viper.GetString("plugins.dockerbuilder.registry_password"),
			Email:    viper.GetString("plugins.dockerbuilder.registry_user_email"),
		},
		BuildArgs: []plugins.Arg{},
	}

	suite.agent.Events <- agent.NewEvent(dockerBuildEvent, nil)

	e = suite.agent.GetTestEvent("plugins.DockerBuild:status", 60)
	payload := e.Payload.(plugins.DockerBuild)
	assert.Equal(suite.T(), string(payload.Action), string(plugins.Status))
	assert.Equal(suite.T(), string(payload.State), string(plugins.Fetching))

	e = suite.agent.GetTestEvent("plugins.DockerBuild:status", 600)
	payload = e.Payload.(plugins.DockerBuild)
	assert.Equal(suite.T(), string(payload.Action), string(plugins.Status))
	assert.Equal(suite.T(), string(payload.State), string(plugins.Complete))
}

func TestDockerBuilder(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
