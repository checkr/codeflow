package webhooks

import (
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"

	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins/webhooks/github"
)

type Webhook interface {
	Register(router *mux.Router, events chan agent.Event)
}

func init() {
	agent.RegisterPlugin("webhooks", func() agent.Plugin { return NewWebhooks() })
}

type Webhooks struct {
	ServiceAddress string `mapstructure:"service_address"`

	Github *github.GithubWebhook
}

func NewWebhooks() *Webhooks {
	return &Webhooks{}
}

func (wb *Webhooks) SampleConfig() string {
	return `
  ## Address and port to host Webhook listener on
	service_address: ":1619"
	github:
		path: "/github"
	`
}

func (wb *Webhooks) Description() string {
	return "A Webhooks Event collector"
}

func (wb *Webhooks) Listen(events chan agent.Event) {
	r := mux.NewRouter()

	for _, webhook := range wb.AvailableWebhooks() {
		webhook.Register(r, events)
	}

	err := http.ListenAndServe(fmt.Sprintf("%s", wb.ServiceAddress), r)
	if err != nil {
		log.Printf("Error starting server: %v", err)
	}
}

// Looks for fields which implement Webhook interface
func (wb *Webhooks) AvailableWebhooks() []Webhook {
	webhooks := make([]Webhook, 0)
	s := reflect.ValueOf(wb).Elem()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if !f.CanInterface() {
			continue
		}

		if wbPlugin, ok := f.Interface().(Webhook); ok {
			if !reflect.ValueOf(wbPlugin).IsNil() {
				webhooks = append(webhooks, wbPlugin)
			}
		}
	}

	return webhooks
}

func (wb *Webhooks) Start(events chan agent.Event) error {
	go wb.Listen(events)
	log.Printf("Started the webhooks service on %s\n", wb.ServiceAddress)
	return nil
}

func (rb *Webhooks) Stop() {
	log.Println("Stopping the Webhooks service")
}

func (rb *Webhooks) Subscribe() []string {
	return []string{}
}

func (rb *Webhooks) Process(e agent.Event) error {
	log.Printf("Process Webhook event: %s", e.Name)
	return nil
}
