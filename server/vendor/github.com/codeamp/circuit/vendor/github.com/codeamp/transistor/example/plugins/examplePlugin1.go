package plugins

import (
	log "github.com/codeamp/logger"
	"github.com/codeamp/transistor"
)

type ExamplePlugin1 struct {
	events chan transistor.Event
	Hello  string `mapstructure:"hello"`
}

func init() {
	transistor.RegisterPlugin("examplePlugin1", func() transistor.Plugin {
		return &ExamplePlugin1{}
	})
}

func (x *ExamplePlugin1) Start(e chan transistor.Event) error {
	log.Info("starting ExamplePlugin")

	event := Hello{
		Action:  "examplePlugin2",
		Message: "Hello World from ExamplePlugin1",
	}

	e <- transistor.NewEvent(event, nil)

	return nil
}

func (x *ExamplePlugin1) Stop() {
	log.Info("stopping ExamplePlugin")
}

func (x *ExamplePlugin1) Subscribe() []string {
	return []string{
		"plugins.Hello:examplePlugin2",
	}
}

func (x *ExamplePlugin1) Process(e transistor.Event) error {
	if e.Name == "plugins.Hello:examplePlugin2" {
		hello := e.Payload.(Hello)
		log.Info("ExamplePlugin1 received a message:", hello)
	}
	return nil
}
