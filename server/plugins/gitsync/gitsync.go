package gitsync

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
	log "github.com/codeamp/logger"
	"github.com/spf13/viper"
)

type GitSync struct {
	events chan agent.Event
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

func (x *GitSync) commits(project plugins.Project, git plugins.Git) ([]plugins.GitCommit, error) {
	var err error
	var gitClone []byte
	var gitPull []byte
	var gitCheckout []byte

	idRsaPath := fmt.Sprintf("%s/%s_id_rsa", viper.GetString("plugins.gitsync.workdir"), project.Repository)
	idRsa := fmt.Sprintf("GIT_SSH_COMMAND=\"ssh -i %s -F /dev/null\"", idRsaPath)
	repoPath := fmt.Sprintf("%s/%s_%s", viper.GetString("plugins.gitsync.workdir"), project.Repository, git.Branch)

	if _, err = os.Stat(repoPath); err != nil {
		if os.IsNotExist(err) {
			mkdir, err := exec.Command("mkdir", "-p", repoPath).CombinedOutput()
			if err != nil {
				log.Debug(err)
				return nil, err
			}
			log.Info(string(mkdir))
		}
	}

	err = ioutil.WriteFile(idRsaPath, []byte(git.RsaPrivateKey), 0600)
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	if _, err = os.Stat(fmt.Sprintf("%s/.git", repoPath)); err != nil {
		if os.IsNotExist(err) {
			run := exec.Command("git", "clone", "-b", git.Branch, "--single-branch", git.Url, repoPath)
			run.Env = []string{idRsa}
			gitClone, err = run.CombinedOutput()
			if err != nil {
				log.Debug(err)
				return nil, err
			}
			log.Info(string(gitClone))
		}
	} else {
		run := exec.Command("git", "-C", repoPath, "pull", "origin", git.Branch)
		run.Env = []string{idRsa}
		gitPull, err = run.CombinedOutput()
		if err != nil {
			log.Debug(err)
			return nil, err
		}
		log.Info(string(gitPull))
	}

	run := exec.Command("git", "-C", repoPath, "checkout", git.Branch)
	run.Env = []string{idRsa}
	gitCheckout, err = run.CombinedOutput()
	if err != nil {
		log.Debug(err)
		return nil, err
	}
	log.Info(string(gitCheckout))

	gitLog, err := exec.Command("git", "log", "--date=iso-strict", "-n", "50", `--pretty=format:{###hash###:###%H###,###parentHash###:###%P###,###message###:###%s###,###user###:###%cN###,###created###:###%cd###},`).Output()

	if err != nil {
		log.Debug(err)
		return nil, err
	}

	gitLogString := string(gitLog)
	gitLogString = strings.Replace(gitLogString, "\"", "'", -1)
	gitLogString = strings.Replace(gitLogString, "{###", "{\"", -1)
	gitLogString = strings.Replace(gitLogString, "###}", "\"}", -1)
	gitLogString = strings.Replace(gitLogString, "###:###", "\":\"", -1)
	gitLogString = strings.Replace(gitLogString, "###,###", "\",\"", -1)

	gitLogJson := fmt.Sprintf("[%s]", gitLogString[:len(gitLogString)-1])

	var commits []plugins.GitCommit

	err = json.Unmarshal([]byte(gitLogJson), &commits)
	if err != nil {
		log.Debug(err)
		return nil, err
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
