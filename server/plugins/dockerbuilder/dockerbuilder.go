package gitsync

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
	log "github.com/codeamp/logger"
	"github.com/extemporalgenome/slug"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/spf13/viper"
)

type DockerBuilder struct {
	events chan agent.Event
	Socket string
}

func init() {
	agent.RegisterPlugin("dockerbuilder", func() agent.Plugin {
		return &DockerBuilder{Socket: "unix:///var/run/docker.sock"}
	})
}

func (x *DockerBuilder) Description() string {
	return "Clone git repository and build a docker image"
}

func (x *DockerBuilder) SampleConfig() string {
	return ` `
}

func (x *DockerBuilder) Start(e chan agent.Event) error {
	x.events = e
	log.Info("Started DockerBuilder")

	return nil
}

func (x *DockerBuilder) Stop() {
	log.Println("Stopping DockerBuilder")
}

func (x *DockerBuilder) Subscribe() []string {
	return []string{
		"plugins.DockerBuild:create",
	}
}

func (x *DockerBuilder) git(env []string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)

	log.InfoWithFields("executing command", log.Fields{
		"path": cmd.Path,
		"args": strings.Join(cmd.Args, " "),
	})

	cmd.Env = env

	out, err := cmd.CombinedOutput()
	if err != nil {
		if ee, ok := err.(*exec.Error); ok {
			if ee.Err == exec.ErrNotFound {
				return nil, errors.New("Git executable not found in $PATH")
			}
		}

		return nil, errors.New(string(bytes.TrimSpace(out)))
	}

	return out, nil
}

func (x *DockerBuilder) bootstrap(repoPath string, imagePath string, event plugins.DockerBuild) error {
	var err error
	var output []byte

	idRsaPath := fmt.Sprintf("%s/%s_id_rsa", event.Git.Workdir, event.Project.Repository)
	idRsa := fmt.Sprintf("GIT_SSH_COMMAND=ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -i %s -F /dev/null", idRsaPath)

	// Git Env
	env := os.Environ()
	env = append(env, idRsa)

	log.Debug(repoPath)
	_, err = exec.Command("mkdir", "-p", filepath.Dir(repoPath)).CombinedOutput()
	if err != nil {
		return err
	}

	if _, err := os.Stat(idRsaPath); os.IsNotExist(err) {
		log.InfoWithFields("creating repository id_rsa", log.Fields{
			"path": idRsaPath,
		})

		err := ioutil.WriteFile(idRsaPath, []byte(event.Git.RsaPrivateKey), 0600)
		if err != nil {
			log.Debug(err)
			return err
		}
	}

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		log.InfoWithFields("cloning repository", log.Fields{
			"path": repoPath,
		})

		output, err := x.git(env, "clone", event.Git.Url, repoPath)
		if err != nil {
			log.Debug(err)
			return err
		}
		log.Info(string(output))
	}

	output, err = x.git(env, "-C", repoPath, "pull", "origin", event.Git.Branch)
	if err != nil {
		log.Debug(err)
		return err
	}
	log.Info(string(output))

	output, err = x.git(env, "-C", repoPath, "checkout", event.Git.Branch)
	if err != nil {
		log.Debug(err)
		return err
	}
	log.Info(string(output))

	return nil
}

func (x *DockerBuilder) build(repoPath string, nameTag string, event plugins.DockerBuild, dockerBuildOut io.Writer) error {
	gitArchive := exec.Command("git", "archive", event.Feature.Hash)
	gitArchive.Dir = repoPath

	gitArchiveOut, err := gitArchive.StdoutPipe()
	if err != nil {
		log.Debug(err)
		return err
	}

	gitArchiveErr, err := gitArchive.StderrPipe()
	if err != nil {
		log.Debug(err)
		return err
	}

	err = gitArchive.Start()
	if err != nil {
		log.Fatal(err)
		return err
	}

	dockerBuildIn := bytes.NewBuffer(nil)

	go func() {
		io.Copy(os.Stderr, gitArchiveErr)
	}()

	io.Copy(dockerBuildIn, gitArchiveOut)

	err = gitArchive.Wait()
	if err != nil {
		log.Debug(err)
		return err
	}

	var buildArgs []docker.BuildArg
	for _, arg := range event.BuildArgs {
		ba := docker.BuildArg{
			Name:  arg.Key,
			Value: arg.Value,
		}
		buildArgs = append(buildArgs, ba)
	}

	buildOptions := docker.BuildImageOptions{
		Dockerfile:   "Dockerfile",
		Name:         nameTag,
		OutputStream: dockerBuildOut,
		InputStream:  dockerBuildIn,
		BuildArgs:    buildArgs,
	}

	dockerClient, err := docker.NewClient(x.Socket)
	if err != nil {
		return err
	}

	err = dockerClient.BuildImage(buildOptions)
	if err != nil {
		return err
	}

	return nil
}

func (x *DockerBuilder) push(repoPath string, imagePath string, event plugins.DockerBuild, buildlog io.Writer) error {
	var err error

	buildlog.Write([]byte(fmt.Sprintf("Pushing %s\n", imagePath)))

	dockerClient, err := docker.NewClient(x.Socket)

	imagePathSplit := strings.Split(imagePath, ":")

	tag_latest := "latest"

	if viper.GetString("environment") != "production" {
		tag_latest = fmt.Sprintf("%s.%s", "latest", viper.GetString("environment"))
	}

	err = dockerClient.PushImage(docker.PushImageOptions{
		Name:         imagePathSplit[0],
		Tag:          imagePathSplit[1],
		OutputStream: buildlog,
	}, docker.AuthConfiguration{
		Username: event.Registry.Username,
		Password: event.Registry.Password,
		Email:    event.Registry.Email,
	})
	if err != nil {
		return err
	}

	tagOptions := docker.TagImageOptions{
		Repo:  imagePathSplit[0],
		Tag:   tag_latest,
		Force: true,
	}

	if err = dockerClient.TagImage(imagePath, tagOptions); err != nil {
		return err
	}

	err = dockerClient.PushImage(docker.PushImageOptions{
		Name:         imagePathSplit[0],
		Tag:          tag_latest,
		OutputStream: buildlog,
	}, docker.AuthConfiguration{
		Username: event.Registry.Username,
		Password: event.Registry.Password,
		Email:    event.Registry.Email,
	})
	if err != nil {
		return err
	}

	return nil
}

func (x *DockerBuilder) Process(e agent.Event) error {
	log.InfoWithFields("Process DockerBuilder event", log.Fields{
		"event": e.Name,
	})

	var err error

	event := e.Payload.(plugins.DockerBuild)
	event.Action = plugins.Status
	event.State = plugins.Fetching
	event.StateMessage = ""
	x.events <- e.NewEvent(event, nil)

	repoPath := fmt.Sprintf("%s/%s_%s", event.Git.Workdir, event.Project.Repository, event.Git.Branch)
	imagePath := fmt.Sprintf("%s/%s/%s:%s.%s", event.Registry.Host, event.Registry.Org, slug.Slug(event.Project.Repository), event.Feature.Hash, viper.GetString("environment"))

	buildlog := bytes.NewBuffer(nil)
	//buildlog := io.MultiWriter(buf, os.Stdout)

	err = x.bootstrap(repoPath, imagePath, event)
	if err != nil {
		log.Debug(err)
		event.State = plugins.Failed
		event.StateMessage = fmt.Sprintf("%v (Action: %v, Step: bootstrap)", err.Error(), event.State)
		event := e.NewEvent(event, err)
		x.events <- event
		return err
	}

	err = x.build(repoPath, imagePath, event, buildlog)
	if err != nil {
		log.Debug(err)
		event.State = plugins.Failed
		event.StateMessage = fmt.Sprintf("%v (Action: %v, Step: build)", err.Error(), event.State)
		event.BuildLog = buildlog.String()
		event := e.NewEvent(event, err)
		x.events <- event
		return err
	}

	err = x.push(repoPath, imagePath, event, buildlog)
	if err != nil {
		log.Debug(err)
		event.State = plugins.Failed
		event.StateMessage = fmt.Sprintf("%v (Action: %v, Step: push)", err.Error(), event.State)
		event.BuildLog = buildlog.String()
		event := e.NewEvent(event, err)
		x.events <- event
		return err
	}

	event.State = plugins.Complete
	event.StateMessage = ""
	event.BuildLog = buildlog.String()
	event.Image = imagePath
	x.events <- e.NewEvent(event, nil)
	return nil
}
