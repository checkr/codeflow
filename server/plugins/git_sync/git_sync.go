package git_sync

import (
	"fmt"
	"log"

	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
)

type GitSync struct {
	events chan agent.Event
}

func init() {
	agent.RegisterPlugin("git_sync", func() agent.Plugin {
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
	log.Println("Started GitSync")

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

func (x *GitSync) Process(e agent.Event) error {
	log.Printf("Process GitSync event: %s", e.Name)
	var err error

	gitSyncEvent := e.Payload.(plugins.GitSync)
	gitSyncEvent.Action = plugins.Status
	gitSyncEvent.State = plugins.Fetching
	gitSyncEvent.StateMessage = ""

	commits, err := plugins.GitCommits(gitSyncEvent.HeadHash, gitSyncEvent.Project, gitSyncEvent.Git)
	if err != nil {
		log.Println(err)
		gitSyncEvent.State = plugins.Failed
		gitSyncEvent.StateMessage = fmt.Sprintf("%v (Action: %v)", err.Error(), gitSyncEvent.State)
		event := e.NewEvent(gitSyncEvent, err)
		x.events <- event
		return err
	}

	for i, _ := range commits {
		c := commits[i]
		x.events <- e.NewEvent(c, nil)
	}
	return nil
}
