package gitsync

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/codeamp/circuit/plugins"
	log "github.com/codeamp/logger"
	"github.com/codeamp/transistor"
	"github.com/spf13/viper"
)

type GitSync struct {
	events chan transistor.Event
}

func init() {
	transistor.RegisterPlugin("gitsync", func() transistor.Plugin {
		return &GitSync{}
	})
}

func (x *GitSync) Start(e chan transistor.Event) error {
	x.events = e
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

func (x *GitSync) commits(project plugins.Project, git plugins.Git) ([]plugins.GitCommit, error) {
	var err error

	idRsaPath := fmt.Sprintf("%s/%s_id_rsa", viper.GetString("plugins.gitsync.workdir"), project.Repository)
	idRsa := fmt.Sprintf("GIT_SSH_COMMAND=\"ssh -i %s -F /dev/null\"", idRsaPath)
	repoPath := fmt.Sprintf("%s/%s_%s", viper.GetString("plugins.gitsync.workdir"), project.Repository, git.Branch)

	err = ioutil.WriteFile(idRsaPath, []byte(git.RsaPrivateKey), 0600)
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	gitClone := exec.Command("")

	if _, err = os.Stat(repoPath); err != nil {
		if os.IsNotExist(err) {
			gitClone = exec.Command("git", "clone", "-b", git.Branch, "--single-branch", git.Url, repoPath)
		}
	} else {
		gitClone = exec.Command("git", "-C", repoPath, "pull", "origin", git.Branch)
	}

	gitCloneOut, err := gitClone.StdoutPipe()
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	gitCloneErr, err := gitClone.StderrPipe()
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	gitCheckout := exec.Command("git", "-C", repoPath, "checkout", git.Branch)

	gitCheckoutOut, err := gitCheckout.StdoutPipe()
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	gitCheckoutErr, err := gitCheckout.StderrPipe()
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	go func() {
		io.Copy(os.Stdout, gitCloneOut)
	}()

	go func() {
		io.Copy(os.Stderr, gitCloneErr)
	}()

	go func() {
		io.Copy(os.Stdout, gitCheckoutOut)
	}()

	go func() {
		io.Copy(os.Stderr, gitCheckoutErr)
	}()

	gitClone.Env = []string{idRsa}

	err = gitClone.Start()
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	err = gitClone.Wait()
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	err = gitCheckout.Start()
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	err = gitCheckout.Wait()
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	gitLog, err := exec.Command("git", "log", "--date=iso-strict", "-n", "50", `--pretty=format:{
		"hash": "%H",
		"parentHash": "%P",
		"message": "%s",
		"user": "%cN",
		"created": "%cd"
	},`).Output()

	if err != nil {
		log.Debug(err)
		return nil, err
	}

	gitLogJson := string(gitLog)
	gitLogJson = fmt.Sprintf("[%s]", gitLogJson[:len(gitLogJson)-1])

	var commits []plugins.GitCommit
	err = json.Unmarshal([]byte(gitLogJson), &commits)
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	return commits, nil
}

func (x *GitSync) Process(e transistor.Event) error {
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
