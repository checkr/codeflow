package docker_build

import (
	"bytes"
	"fmt"

	"github.com/checkr/codeflow/server/plugins"
	"github.com/extemporalgenome/slug"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/spf13/viper"
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
	repoPath := fmt.Sprintf("%s/%s_%s", build.Git.Workdir, build.Project.Repository, build.Git.Protocol)
	name := fmt.Sprintf("%s/%s/%s:%s.%s", build.Registry.Host, build.Registry.Org, slug.Slug(build.Project.Repository), build.Feature.Hash, viper.GetString("environment"))

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

func (b *DockerBuilder) push(build *plugins.DockerBuild) error {
	var err error
	name := fmt.Sprintf("%s/%s/%s", build.Registry.Host, build.Registry.Org, slug.Slug(build.Project.Repository))
	tag_hash := fmt.Sprintf("%s.%s", build.Feature.Hash, viper.GetString("environment"))
	full_name := fmt.Sprintf("%s:%s", name, tag_hash)
	tag_latest := "latest"

	if viper.GetString("environment") != "production" {
		tag_latest = fmt.Sprintf("%s.%s", "latest", viper.GetString("environment"))
	}

	b.outputBuffer.Write([]byte(fmt.Sprintf("Pushing %s:%s.%s...", name, build.Feature.Hash, viper.GetString("environment"))))

	err = b.dockerClient.PushImage(docker.PushImageOptions{
		Name:         name,
		Tag:          tag_hash,
		OutputStream: b.outputBuffer,
	}, docker.AuthConfiguration{
		Username: build.Registry.Username,
		Password: build.Registry.Password,
		Email:    build.Registry.Email,
	})
	if err != nil {
		return err
	}

	tagOptions := docker.TagImageOptions{
		Repo:  name,
		Tag:   tag_latest,
		Force: true,
	}
	if err = b.dockerClient.TagImage(full_name, tagOptions); err != nil {
		return err
	}

	err = b.dockerClient.PushImage(docker.PushImageOptions{
		Name:         name,
		Tag:          tag_latest,
		OutputStream: b.outputBuffer,
	}, docker.AuthConfiguration{
		Username: build.Registry.Username,
		Password: build.Registry.Password,
		Email:    build.Registry.Email,
	})
	if err != nil {
		return err
	}

	build.Image = full_name

	return nil
}

func (b *DockerBuilder) cleanup(build *plugins.DockerBuild) error {
	return nil
}
