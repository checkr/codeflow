[![GoDoc](https://godoc.org/github.com/ashwanthkumar/slack-go-webhook?status.svg)](https://godoc.org/github.com/ashwanthkumar/slack-go-webhook)

# slack-go-webhook

Go Lang library to send messages to Slack via Incoming Webhooks.

## Usage
```go
package main

import "github.com/ashwanthkumar/slack-go-webhook"
import "fmt"

func main() {
    webhookUrl := "https://hooks.slack.com/services/foo/bar/baz"

    attachment1 := slack.Attachment {}
    attachment1.AddField(slack.Field { Title: "Author", Value: "Ashwanth Kumar" }).AddField(slack.Field { Title: "Status", Value: "Completed" })
    payload := slack.Payload {
      Text: "Hello from <https://github.com/ashwanthkumar/slack-go-webhook|slack-go-webhook>, a Go-Lang library to send slack webhook messages.\n<https://golangschool.com/wp-content/uploads/golang-teach.jpg|golang-img>",
      Username: "robot",
      Channel: "#general",
      IconEmoji: ":monkey_face:",
      Attachments: []slack.Attachment{attachment1},
    }
    err := slack.Send(webhookUrl, "", payload)
    if len(err) > 0 {
      fmt.Printf("error: %s\n", err)
    }
}
```

## License
Licensed under the Apache License, Version 2.0: http://www.apache.org/licenses/LICENSE-2.0

