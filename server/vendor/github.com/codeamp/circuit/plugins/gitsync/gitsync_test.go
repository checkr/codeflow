package gitsync_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/codeamp/circuit/plugins"
	"github.com/codeamp/transistor"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	transistor *transistor.Transistor
}

var viperConfig = []byte(`
plugins:
  gitsync:
    workers: 1
    workdir: /tmp/gitsync
`)

func (suite *TestSuite) SetupSuite() {
	viper.SetConfigType("YAML")
	viper.ReadConfig(bytes.NewBuffer(viperConfig))
	viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")

	config := transistor.Config{
		Plugins:        viper.GetStringMap("plugins"),
		EnabledPlugins: []string{"gitsync"},
	}

	ag, _ := transistor.NewTestTransistor(config)
	suite.transistor = ag
	go suite.transistor.Run()
}

func (suite *TestSuite) TearDownSuite() {
	suite.transistor.Stop()
}

func (suite *TestSuite) TestGitSync() {
	var e transistor.Event

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

	suite.transistor.Events <- transistor.NewEvent(gitSync, nil)

	created := time.Now()
	for i := 0; i < 5; i++ {
		e = suite.transistor.GetTestEvent("plugins.GitCommit", 60)
		payload := e.Payload.(plugins.GitCommit)
		assert.Equal(suite.T(), payload.Repository, gitSync.Project.Repository)
		assert.Equal(suite.T(), payload.Ref, fmt.Sprintf("refs/heads/%s", gitSync.Git.Branch))
		assert.True(suite.T(), payload.Created.Before(created), "Commit created time is older than previous commit")
	}
}

func TestGitSync(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
