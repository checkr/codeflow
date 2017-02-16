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
	case "plugins.DockerDeploy:status":
		payload := e.Payload.(plugins.DockerDeploy)

		if payload.State != plugins.Complete && payload.State != plugins.Failed {
			return nil
		}

		for _, channel := range payload.Project.NotifyChannels {
			project := payload.Project.Slug
			message := payload.Release.HeadFeature.Message
			author := payload.Release.User

			repository := payload.Project.Repository
			tail := payload.Release.TailFeature.Hash
			head := payload.Release.HeadFeature.Hash

			text := fmt.Sprintf(
				"%s deployed <https://github.com/%s/compare/%s...%s|%s...%s> to <%s/projects/%s/deploy/|%s>",
				author, repository, tail, head, tail[0:6], head[0:6], viper.GetString("plugins.codeflow.dashboard_url"), project, project,
			)
			feature_attachment := slack_webhook.Attachment{Text: &message}

			var result_color, result_text, result_emoji string
			if payload.State == plugins.Failed {
				result_color = "#FF0000"
				result_text = "FAILED"
				result_emoji = ":ambulance:"
			} else {
				result_color = "#008000"
				result_text = "SUCCESS"
				result_emoji = ":rocket:"
			}
			result_attachment := slack_webhook.Attachment{Color: &result_color, Text: &result_text}

			slackPayload := slack_webhook.Payload{
				Text:        text,
				Username:    "Codeflow",
				Channel:     channel,
				IconEmoji:   result_emoji,
				Attachments: []slack_webhook.Attachment{feature_attachment, result_attachment},
			}
			err := slack_webhook.Send(webhookUrl, "", slackPayload)
			if len(err) > 0 {
				log.Println("error: %s\n", err)
			}
		}
	}

	return nil
}
