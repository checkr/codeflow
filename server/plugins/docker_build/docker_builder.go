package docker_build

import (
	"bytes"
	"fmt"
	"os"

	"github.com/checkr/codeflow/server/plugins"
	docker "github.com/fsouza/go-dockerclient"
	git "github.com/libgit2/git2go"
)

type DockerBuilder struct {
	dockerClient  *docker.Client
	buildPath     string
	rsaPrivateKey string
	rsaPublicKey  string
	outputBuffer  *bytes.Buffer
}

func NewDockerBuilder(
	dockerClient *docker.Client,
	buildPath string,
	rsaPrivateKey string,
	rsaPublicKey string,
	outputBuffer *bytes.Buffer,
) *DockerBuilder {
	return &DockerBuilder{
		dockerClient:  dockerClient,
		buildPath:     buildPath,
		rsaPrivateKey: rsaPrivateKey,
		rsaPublicKey:  rsaPublicKey,
		outputBuffer:  outputBuffer,
	}
}

func (b *DockerBuilder) fetchCode(build *plugins.DockerBuild) error {
	repoPath := fmt.Sprintf("%s/%s", b.buildPath, build.Project.Repository)

	repo, err := b.findOrClone(repoPath, build.Git.SshUrl)
	if err != nil {
		return err
	}

	remote, err := repo.Remotes.Lookup("origin")
	if err != nil {
		return err
	}

	err = remote.Fetch([]string{}, &git.FetchOptions{
		RemoteCallbacks: git.RemoteCallbacks{
			CredentialsCallback:      b.credentialsCallback,
			CertificateCheckCallback: b.certificateCheckCallback,
		},
	}, "")
	if err != nil {
		return err
	}

	oid, err := git.NewOid(build.Feature.Hash)
	if err != nil {
		return err
	}

	commit, err := repo.LookupCommit(oid)
	if err != nil {
		return err
	}

	tree, err := commit.Tree()
	if err != nil {
		return err
	}

	err = repo.CheckoutTree(tree, &git.CheckoutOpts{Strategy: git.CheckoutForce})
	if err != nil {
		return err
	}

	err = repo.SetHeadDetached(oid)
	if err != nil {
		return err
	}

	return nil
}

func (b *DockerBuilder) build(build *plugins.DockerBuild) error {
	repoPath := fmt.Sprintf("%s/%s", b.buildPath, build.Project.Repository)
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

func (b *DockerBuilder) findOrClone(path string, cloneUrl string) (*git.Repository, error) {
	var repo *git.Repository
	var err error

	if _, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			cloneOptions := &git.CloneOptions{}
			cloneOptions.FetchOptions = &git.FetchOptions{
				RemoteCallbacks: git.RemoteCallbacks{
					CredentialsCallback:      b.credentialsCallback,
					CertificateCheckCallback: b.certificateCheckCallback,
				},
			}
			repo, err = git.Clone(cloneUrl, path, cloneOptions)
		} else {
			return &git.Repository{}, err
		}
	} else {
		repo, err = git.OpenRepository(path)
	}

	return repo, err
}

func (b *DockerBuilder) credentialsCallback(url string, username string, allowedTypes git.CredType) (git.ErrorCode, *git.Cred) {
	id_rsa_priv := b.rsaPrivateKey
	id_rsa_pub := b.rsaPublicKey
	ret, cred := git.NewCredSshKeyFromMemory("git", id_rsa_pub, id_rsa_priv, "")
	return git.ErrorCode(ret), &cred
}

// Made this one just return 0 during troubleshooting...
func (b *DockerBuilder) certificateCheckCallback(cert *git.Certificate, valid bool, hostname string) git.ErrorCode {
	return 0
}
