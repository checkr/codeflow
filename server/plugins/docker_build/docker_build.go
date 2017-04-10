package docker_build

import (
	"bytes"
	"fmt"
	"log"

	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/spf13/viper"
)

type DockerBuild struct {
	events chan agent.Event
}

func init() {
	agent.RegisterPlugin("docker_build", func() agent.Plugin {
		return &DockerBuild{}
	})
}

func (b *DockerBuild) Description() string {
	return "Clone git repository and build a docker image"
}

func (b *DockerBuild) SampleConfig() string {
	return ` `
}

func (b *DockerBuild) Start(events chan agent.Event) error {
	b.events = events
	log.Println("Started Git Docker Build")
	return nil
}

func (b *DockerBuild) Stop() {
	log.Println("Stopping Git Docker Build")
}

func (b *DockerBuild) Subscribe() []string {
	return []string{
		"plugins.DockerBuild:create",
	}
}

func (b *DockerBuild) Process(e agent.Event) error {
	log.Printf("Process DockerBuild event: %s", e.Name)

	var err error
	var event agent.Event

	build := e.Payload.(plugins.DockerBuild)
	build.Action = plugins.Status
	build.State = plugins.Running
	build.BuildLog = ""
	build.BuildError = ""

	dockerHost := viper.GetString("plugins.docker_build.docker_host")
	rsaPrivateKey := build.Git.RsaPrivateKey
	rsaPublicKey := build.Git.RsaPublicKey
	outputBuffer := bytes.NewBuffer(nil)
	dockerClient, err := docker.NewClient(dockerHost)
	if err != nil {
		build.BuildError = fmt.Sprintf("%v (Action: %v)", err.Error(), build.State)
		build.State = plugins.Failed
		event := e.NewEvent(build, err)
		b.events <- event
		log.Println(err)
		return err
	}

	build.State = plugins.Fetching
	dockerBuilder := NewDockerBuilder(dockerClient, rsaPrivateKey, rsaPublicKey, outputBuffer)
	err = dockerBuilder.fetchCode(&build)
	if err != nil {
		build.BuildLog = outputBuffer.String()
		build.BuildError = fmt.Sprintf("%v (Action: %v)", err.Error(), build.State)
		build.State = plugins.Failed
		event := e.NewEvent(build, err)
		b.events <- event
		log.Println(err)
		return err
	}

	build.State = plugins.Building
	event = e.NewEvent(build, nil)
	b.events <- event

	if err = dockerBuilder.build(&build); err != nil {
		build.BuildLog = outputBuffer.String()
		build.BuildError = fmt.Sprintf("%v (Action: %v)", err.Error(), build.State)
		build.State = plugins.Failed
		event := e.NewEvent(build, err)
		b.events <- event
		log.Println(err)
		return err
	}

	build.State = plugins.Pushing
	event = e.NewEvent(build, nil)
	b.events <- event

	if err = dockerBuilder.push(&build); err != nil {
		build.BuildLog = outputBuffer.String()
		build.BuildError = fmt.Sprintf("%v (Action: %v)", err.Error(), build.State)
		build.State = plugins.Failed
		event := e.NewEvent(build, err)
		b.events <- event
		log.Println(err)
		return err
	}

	build.State = plugins.Complete
	build.BuildLog = outputBuffer.String()
	event = e.NewEvent(build, nil)
	b.events <- event

	if err = dockerBuilder.cleanup(&build); err != nil {
		build.BuildLog = outputBuffer.String()
		build.BuildError = fmt.Sprintf("%v (Action: %v)", err.Error(), build.State)
		build.State = plugins.Failed
		event := e.NewEvent(build, err)
		b.events <- event
		log.Println(err)
		return err
	}

	return nil
}
