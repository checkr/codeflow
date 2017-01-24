package route53

import (
	"log"

	"github.com/checkr/codeflow/server/agent"
)

type Route53 struct {
	events chan agent.Event
}

func init() {
	agent.RegisterPlugin("route53", func() agent.Plugin {
		return &Route53{}
	})
}

func (x *Route53) Description() string {
	return "Set Route53 DNS for Kubernetes services"
}

func (x *Route53) SampleConfig() string {
	return ` `
}

func (x *Route53) Start(e chan agent.Event) error {
	x.events = e
	log.Println("Started Route53")
	return nil
}

func (x *Route53) Stop() {
	log.Println("Stopping Route53")
}

func (x *Route53) Subscribe() []string {
	return []string{
		"plugins.LoadBalancer:info",
	}
}

func (x *Route53) Process(e agent.Event) error {

	switch e.Name {
	case "plugins.LoadBalancer:status":
		log.Printf("Process Route53 event: %s", e.Name)
		x.updateRoute53(e)
	}
	return nil
}

func (x *Route53) updateRoute53(e agent.Event) error {
	return nil
}
