package docker_build

import (
	"bytes"
	"fmt"

	"github.com/checkr/codeflow/server/plugins"
	docker "github.com/fsouza/go-dockerclient"
)

type DockerBuilder struct {
	dockerClient  *docker.Client
	rsaPrivateKey string
	rsaPublicKey  string
	outputBuffer  *bytes.Buffer
}

func NewDockerBuilder(
	dockerClient *docker.Client,
	rsaPrivateKey string,
	rsaPublicKey string,
	outputBuffer *bytes.Buffer,
) *DockerBuilder {
	return &DockerBuilder{
		dockerClient:  dockerClient,
		rsaPrivateKey: rsaPrivateKey,
		rsaPublicKey:  rsaPublicKey,
		outputBuffer:  outputBuffer,
	}
}

func (b *DockerBuilder) fetchCode(build *plugins.DockerBuild) error {
	_, err := plugins.GitCheckoutCommit(build.Feature.Hash, build.Project, build.Git)
	if err != nil {
		return err
	}

	return nil
}

func (b *DockerBuilder) build(build *plugins.DockerBuild) error {
	repoPath := fmt.Sprintf("%s/%s", build.Git.Workdir, build.Project.Repository)
	name := fmt.Sprintf("%s/%s:%s.%s", build.Registry.Host, build.Project.Repository, build.Feature.Hash, "codeflow")

	var buildArgs []docker.BuildArg
	for _, arg := range build.BuildArgs {
		ba := docker.BuildArg{
			Name:  arg.Key,
			Value: arg.Value,
		}
		buildArgs = append(buildArgs, ba)
	}

	buildOptions := docker.BuildImageOptions{
		Dockerfile:   "Dockerfile",
		Name:         name,
		OutputStream: b.outputBuffer,
		BuildArgs:    buildArgs,
		ContextDir:   repoPath,
	}

	if err := b.dockerClient.BuildImage(buildOptions); err != nil {
		return err
	}

	return nil
}

func (b *DockerBuilder) tag(build *plugins.DockerBuild) error {
	name := fmt.Sprintf("%s/%s:%s.%s", build.Registry.Host, build.Project.Repository, build.Feature.Hash, "codeflow")
	tagOptions := docker.TagImageOptions{
		Repo:  fmt.Sprintf("%s/%s", build.Registry.Host, build.Project.Repository),
		Tag:   "latest",
		Force: true,
	}
	if err := b.dockerClient.TagImage(name, tagOptions); err != nil {
		return err
	}
	return nil
}

func (b *DockerBuilder) push(build *plugins.DockerBuild) error {
	name := fmt.Sprintf("%s/%s", build.Registry.Host, build.Project.Repository)

	b.outputBuffer.Write([]byte(fmt.Sprintf("Pushing %s:%s.%s...", name, build.Feature.Hash, "codeflow")))

	err := b.dockerClient.PushImage(docker.PushImageOptions{
		Name:         name,
		Tag:          fmt.Sprintf("%s.%s", build.Feature.Hash, "codeflow"),
		OutputStream: b.outputBuffer,
	}, docker.AuthConfiguration{
		Username: build.Registry.Username,
		Password: build.Registry.Password,
		Email:    build.Registry.Email,
	})
	if err != nil {
		return err
	}

	build.Image = fmt.Sprintf("%s/%s:%s.%s", build.Registry.Host, build.Project.Repository, build.Feature.Hash, "codeflow")

	return nil
}

func (b *DockerBuilder) cleanup(build *plugins.DockerBuild) error {
	return nil
}
