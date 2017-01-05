package heartbeat

import (
	"log"
	"time"

	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
)

const ANY = -1

type Heartbeat struct {
	events chan agent.Event
	ticks  []tick
}

type tick struct {
	Month, Day, Weekday  int8
	Hour, Minute, Second int8
	Name                 string
}

func (x *tick) Matches(t time.Time) (ok bool) {
	ok = (x.Month == ANY || x.Month == int8(t.Month())) &&
		(x.Day == ANY || x.Day == int8(t.Day())) &&
		(x.Weekday == ANY || x.Weekday == int8(t.Weekday())) &&
		(x.Hour == ANY || x.Hour == int8(t.Hour())) &&
		(x.Minute == ANY || x.Minute == int8(t.Minute())) &&
		(x.Second == ANY || x.Second == int8(t.Second()))
	return ok
}

func (x *Heartbeat) NewTick(month, day, weekday, hour, minute, second int8, name string) {
	b := tick{month, day, weekday, hour, minute, second, name}
	x.ticks = append(x.ticks, b)
}

func init() {
	agent.RegisterPlugin("heartbeat", func() agent.Plugin {
		return &Heartbeat{}
	})
}

func (x *Heartbeat) Description() string {
	return "Send heartbeat events to subscribed plugins"
}

func (x *Heartbeat) SampleConfig() string {
	return ` `
}

func (x *Heartbeat) heartbeat() error {
	for {
		now := time.Now()
		for _, b := range x.ticks {
			if b.Matches(now) {
				tick := plugins.HeartBeat{Tick: b.Name}
				event := agent.NewEvent(tick, nil)
				x.events <- event
			}
		}
		time.Sleep(time.Second)
	}
}

func (x *Heartbeat) Start(e chan agent.Event) error {
	x.events = e

	x.NewTick(ANY, ANY, ANY, ANY, ANY, 0, "minute")
	x.NewTick(ANY, ANY, ANY, ANY, 0, 0, "hour")
	x.NewTick(ANY, ANY, 0, 0, 0, 0, "day")

	go x.heartbeat()
	log.Println("Started Heartbeat")

	return nil
}

func (x *Heartbeat) Stop() {
	log.Println("Stopping Heartbeat")
}

func (x *Heartbeat) Subscribe() []string {
	return []string{}
}

func (x *Heartbeat) Process(e agent.Event) error {
	return nil
}
