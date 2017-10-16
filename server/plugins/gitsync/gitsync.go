package gitsync

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
	log "github.com/codeamp/logger"
	"github.com/spf13/viper"
)

type GitSync struct {
	events chan agent.Event
	idRsa  string
}

func init() {
	agent.RegisterPlugin("gitsync", func() agent.Plugin {
		return &GitSync{}
	})
}

func (x *GitSync) Description() string {
	return "Sync Git repositories and create new features"
}

func (x *GitSync) SampleConfig() string {
	return ` `
}

func (x *GitSync) Start(e chan agent.Event) error {
	x.events = e
	log.SetLogLevel("debug")
	log.Info("Started GitSync")

	return nil
}

func (x *GitSync) Stop() {
	log.Println("Stopping GitSync")
}

func (x *GitSync) Subscribe() []string {
	return []string{
		"plugins.GitPing",
		"plugins.GitSync:update",
	}
}

func (x *GitSync) git(args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	env := os.Environ()
	env = append(env, x.idRsa)
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

func (x *GitSync) toGitCommit(entry string) (plugins.GitCommit, error) {
	items := strings.Split(entry, "#@#")
	commiterDate, err := time.Parse("2006-01-02T15:04:05-07:00", items[4])

	if err != nil {
		return plugins.GitCommit{}, err
	}

	return plugins.GitCommit{
		User:       items[3],
		Message:    items[2],
		Hash:       items[0],
		ParentHash: items[1],
		Created:    commiterDate,
	}, nil
}

func (x *GitSync) commits(project plugins.Project, git plugins.Git) ([]plugins.GitCommit, error) {
	var err error
	var gitClone []byte
	var gitPull []byte
	var gitCheckout []byte

	idRsaPath := fmt.Sprintf("%s/%s_id_rsa", viper.GetString("plugins.gitsync.workdir"), project.Repository)
	x.idRsa = fmt.Sprintf("GIT_SSH_COMMAND=ssh -i %s -F /dev/null", idRsaPath)
	repoPath := fmt.Sprintf("%s/%s_%s", viper.GetString("plugins.gitsync.workdir"), project.Repository, git.Branch)

	mkdir, err := exec.Command("mkdir", "-p", filepath.Dir(repoPath)).CombinedOutput()
	if err != nil {
		log.Debug(err)
		return nil, err
	}
	log.Info(string(mkdir))

	err = ioutil.WriteFile(idRsaPath, []byte(git.RsaPrivateKey), 0600)
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	if _, err = os.Stat(fmt.Sprintf("%s", repoPath)); err != nil {
		if os.IsNotExist(err) {
			gitClone, err = x.git("clone", git.Url, repoPath)
			if err != nil {
				log.Debug(err)
				return nil, err
			}
			log.Info(string(gitClone))
		}
	} else {
		gitPull, err = x.git("-C", repoPath, "pull", "origin", git.Branch)
		if err != nil {
			log.Debug(err)
			return nil, err
		}
		log.Info(string(gitPull))
	}

	gitCheckout, err = x.git("-C", repoPath, "checkout", git.Branch)
	if err != nil {
		log.Debug(err)
		return nil, err
	}
	log.Info(string(gitCheckout))

	gitLog, err := x.git("-C", repoPath, "log", "--no-merges", "--date=iso-strict", "-n", "50", "--pretty=format:%H#@#%P#@#%s#@#%cN#@#%cd", git.Branch)

	if err != nil {
		log.Debug(err)
		return nil, err
	}

	var commits []plugins.GitCommit

	for _, line := range strings.Split(strings.TrimSuffix(string(gitLog), "\n"), "\n") {
		commit, err := x.toGitCommit(line)
		if err != nil {
			log.Debug(err)
			return nil, err
		}

		commits = append(commits, commit)
	}

	return commits, nil
}

func (x *GitSync) Process(e agent.Event) error {
	log.InfoWithFields("Process GitSync event", log.Fields{
		"event": e.Name,
	})

	var err error

	gitSyncEvent := e.Payload.(plugins.GitSync)
	gitSyncEvent.Action = plugins.Status
	gitSyncEvent.State = plugins.Fetching
	gitSyncEvent.StateMessage = ""
	x.events <- e.NewEvent(gitSyncEvent, nil)

	commits, err := x.commits(gitSyncEvent.Project, gitSyncEvent.Git)
	if err != nil {
		log.Debug(err)
		gitSyncEvent.State = plugins.Failed
		gitSyncEvent.StateMessage = fmt.Sprintf("%v (Action: %v)", err.Error(), gitSyncEvent.State)
		event := e.NewEvent(gitSyncEvent, err)
		x.events <- event
		return err
	}

	for i := range commits {
		c := commits[i]
		c.Repository = gitSyncEvent.Project.Repository
		c.Ref = fmt.Sprintf("refs/heads/%s", gitSyncEvent.Git.Branch)

		if c.Hash == gitSyncEvent.From {
			return nil
		}

		x.events <- e.NewEvent(c, nil)
	}

	gitSyncEvent.State = plugins.Complete
	gitSyncEvent.StateMessage = ""
	x.events <- e.NewEvent(gitSyncEvent, nil)

	return nil
}
