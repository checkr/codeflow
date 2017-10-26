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
			Slug:       "codeamp-circuit",
			Repository: "codeamp/circuit",
		},
		Git: plugins.Git{
			Url:           "https://github.com/codeamp/circuit.git",
			Protocol:      "HTTPS",
			Branch:        "master",
			RsaPrivateKey: "",
			RsaPublicKey:  "",
			Workdir:       viper.GetString("plugins.dockerbuilder.workdir"),
		},
		Feature: plugins.Feature{
			Hash:       "b82f00530a7186d5b03ead5bd3d3600053b71ee7",
			ParentHash: "b5021f702069ac6160fe5f0e9395351a36462c59",
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

	//e = suite.agent.GetTestEvent("plugins.DockerBuild:status", 600)
	//payload = e.Payload.(plugins.DockerBuild)
	//assert.Equal(suite.T(), string(payload.Action), string(plugins.Status))
	//assert.Equal(suite.T(), string(payload.State), string(plugins.Complete))
}

func TestDockerBuilder(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
