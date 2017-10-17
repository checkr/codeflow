package heartbeat

import (
	"time"

	"github.com/codeamp/circuit/plugins"
	log "github.com/codeamp/logger"
	"github.com/codeamp/transistor"
	"github.com/rk/go-cron"
)

type Heartbeat struct {
	events chan transistor.Event
}

func init() {
	transistor.RegisterPlugin("heartbeat", func() transistor.Plugin {
		return &Heartbeat{}
	})
}

func (x *Heartbeat) Start(e chan transistor.Event) error {
	x.events = e

	cron.NewCronJob(cron.ANY, cron.ANY, cron.ANY, cron.ANY, cron.ANY, 0, func(time.Time) {
		event := transistor.NewEvent(plugins.HeartBeat{Tick: "minute"}, nil)
		x.events <- event
	})

	cron.NewCronJob(cron.ANY, cron.ANY, cron.ANY, cron.ANY, 0, 0, func(time.Time) {
		event := transistor.NewEvent(plugins.HeartBeat{Tick: "hour"}, nil)
		x.events <- event
	})

	log.Println("Started Heartbeat")

	return nil
}

func (x *Heartbeat) Stop() {
	log.Println("Stopping Heartbeat")
}

func (x *Heartbeat) Subscribe() []string {
	return []string{}
}

func (x *Heartbeat) Process(e transistor.Event) error {
	return nil
}
