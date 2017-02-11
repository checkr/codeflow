package slack

import (
	"fmt"
	"log"

	slack_webhook "github.com/ashwanthkumar/slack-go-webhook"
	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
	"github.com/spf13/viper"
)

type Slack struct {
}

func init() {
	agent.RegisterPlugin("slack", func() agent.Plugin {
		return &Slack{}
	})
}

func (x *Slack) Description() string {
	return "Send slack events to subscribed plugins"
}

func (x *Slack) SampleConfig() string {
	return ` `
}

func (x *Slack) slack() error {
	return nil
}

func (x *Slack) Start(e chan agent.Event) error {
	go x.slack()
	log.Println("Started Slack")

	return nil
}

func (x *Slack) Stop() {
	log.Println("Stopping Slack")
}

func (x *Slack) Subscribe() []string {
	return []string{
		"plugins.DockerBuild",
		"plugins.DockerDeploy",
		"plugins.Project",
		"plugins.Release",
	}
}

func (x *Slack) Process(e agent.Event) error {
	webhookUrl := viper.GetString("plugins.slack.webhook_url")

	switch e.Name {
	case "plugins.DockerDeploy:create":
		payload := e.Payload.(plugins.DockerDeploy)

		for _, channel := range payload.Project.NotifyChannels {
			project := payload.Project.Slug
			release := payload.Release.HeadFeature.Hash
			message := payload.Release.HeadFeature.Message

			repository := payload.Project.Repository
			tail := payload.Release.TailFeature.Hash
			head := payload.Release.HeadFeature.Hash

			diffUrl := fmt.Sprintf("https://github.com/%s/compare/%s...%s", repository, tail, head)

			attachment1 := slack_webhook.Attachment{}
			attachment1.AddField(slack_webhook.Field{Title: "Commit", Value: message}).AddField(slack_webhook.Field{Title: "Link", Value: diffUrl})

			slackPayload := slack_webhook.Payload{
				Text:        fmt.Sprintf("Deploying %s for %s", release, project),
				Username:    "codeflow-bot",
				Channel:     channel,
				IconEmoji:   ":rocket:",
				Attachments: []slack_webhook.Attachment{attachment1},
			}
			err := slack_webhook.Send(webhookUrl, "", slackPayload)
			if len(err) > 0 {
				log.Println("error: %s\n", err)
			}
		}
	case "plugins.DockerDeploy:status":
		payload := e.Payload.(plugins.DockerDeploy)

		for _, channel := range payload.Project.NotifyChannels {
			project := payload.Project.Slug
			release := payload.Release.HeadFeature.Hash

			msg := fmt.Sprintf("Deploy %s for %s", release, project)
			color := "#00FF00"

			if payload.State == "failed" {
				color = "#FF0000"
				msg = fmt.Sprintf("Deploy %s for %s", release, project)
			}

			attachment1 := slack_webhook.Attachment{Color: &color}
			attachment1.AddField(slack_webhook.Field{Title: "Status", Value: string(payload.State)})

			slackPayload := slack_webhook.Payload{
				Text:        msg,
				Username:    "codeflow-bot",
				Channel:     channel,
				IconEmoji:   ":rocket:",
				Attachments: []slack_webhook.Attachment{attachment1},
			}
			err := slack_webhook.Send(webhookUrl, "", slackPayload)
			if len(err) > 0 {
				log.Println("error: %s\n", err)
			}
		}
	}

	return nil
}
