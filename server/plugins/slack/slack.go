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
			message := payload.Release.HeadFeature.Message
			author  := payload.Release.HeadFeature.User

			repository := payload.Project.Repository
			tail := payload.Release.TailFeature.Hash
			head := payload.Release.HeadFeature.Hash

			msg := fmt.Sprintf(
				"<https://github.com/%s|%s> deploying <https://github.com/%s/compare/%s...%s|%s...%s> to <https://github.com/%s|%s>",
				author, author, repository, tail, head, tail[0:6], head[0:6], repository, project,
			)
			attachment1 := slack_webhook.Attachment{Text: &message}

			slackPayload := slack_webhook.Payload{
				Text:        msg,
				Username:    "Codeflow",
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

		if payload.State != plugins.Complete && payload.State != plugins.Failed {
			return nil
		}

		for _, channel := range payload.Project.NotifyChannels {
			project := payload.Project.Slug
			repository := payload.Project.Repository
			tail := payload.Release.TailFeature.Hash
			head := payload.Release.HeadFeature.Hash

			var msg, color string

			if payload.State == plugins.Failed {
				color = "#FF0000"
				msg = fmt.Sprintf(
					"FAILED | <https://github.com/%s|%s> | <https://github.com/%s/compare/%s...%s|%s...%s>",
					repository, project, repository, tail, head, tail[0:6], head[0:6],
				)
			} else {
				color = "#008000"
				msg = fmt.Sprintf(
					"SUCCEEDED | <https://github.com/%s|%s> | <https://github.com/%s/compare/%s...%s|%s...%s>",
					repository, project, repository, tail, head, tail[0:6], head[0:6],
				)
			}

			attachment1 := slack_webhook.Attachment{Color: &color, Text: &msg}

			slackPayload := slack_webhook.Payload{
				Username:    "Codeflow",
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
