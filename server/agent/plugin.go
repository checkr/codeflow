package agent

import (
	"github.com/jrallison/go-workers"
)

type Plugin interface {
	// SampleConfig returns the default configuration of the Input
	SampleConfig() string

	// Description returns a one-sentence description on the Input
	Description() string

	// Start starts the Plugin service, whatever that may be
	Start(chan Event) error

	// Stop stops the services and closes any necessary channels and connections
	Stop()

	// Subscribe takes in an event message and validates it for Process
	Subscribe() []string

	// Process takes in an event message and tries to process it
	Process(Event) error
}

type RunningPlugin struct {
	Name    string
	Plugin  Plugin
	Work    func(*workers.Msg)
	Enabled bool
	Workers int
}
