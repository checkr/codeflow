package agent

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/jrallison/go-workers"
	"github.com/pborman/uuid"
	"github.com/spf13/viper"
)

// Agent runs codeflow and collects data based on the given config
type Agent struct {
	Events  chan Event
	Plugins []*RunningPlugin

	testRun         bool
	TestSubscribeTo []string
	TestWork        func(*Event, int)
}

// NewAgent returns an Agent struct based off the given Config
func NewAgent() (*Agent, error) {
	if len(viper.GetStringMap("plugins")) == 0 {
		log.Fatalf("Error: no plugins found, did you provide a valid config file?")
	}

	agent := &Agent{}

	if err := agent.LoadPlugins(); err != nil {
		log.Println(err)
	}
	return agent, nil
}

func (a *Agent) LoadPlugins() error {
	var err error

	ep := viper.GetString("run")
	runPlugins := strings.Split(strings.Trim(ep, "[]"), ",")
	for name := range viper.GetStringMap("plugins") {
		if err = a.addPlugin(name); err != nil {
			return fmt.Errorf("Error parsing %s, %s", name, err)
		}
		if ep == "" || SliceContains(name, runPlugins) {
			if err = a.enablePlugin(name); err != nil {
				return fmt.Errorf("Error parsing %s, %s", name, err)
			}
		}
	}

	return nil
}

// Returns a list of strings of the configured plugins.
func (a *Agent) PluginNames() []string {
	var name []string
	for key, _ := range viper.GetStringMap("plugins") {
		name = append(name, key)
	}
	return name
}

func (a *Agent) addPlugin(name string) error {
	if len(a.PluginNames()) > 0 && !SliceContains(name, a.PluginNames()) {
		return nil
	}

	creator, ok := PluginRegistry[name]
	if !ok {
		return fmt.Errorf("Undefined but requested Plugin: %s", name)
	}
	plugin := creator()

	viper.UnmarshalKey(fmt.Sprint("plugins.", name), plugin)

	work := func(message *workers.Msg) {
		e, _ := json.Marshal(message.Args())
		event := Event{}
		json.Unmarshal([]byte(e), &event)
		if err := MapPayload(event.PayloadModel, &event); err != nil {
			event.Error = fmt.Errorf("PayloadModel not found: %s. Did you add it to ApiRegistry?", event.PayloadModel)
		}
		plugin.Process(event)
	}

	rp := &RunningPlugin{
		Name:    name,
		Plugin:  plugin,
		Work:    work,
		Enabled: false,
	}

	a.Plugins = append(a.Plugins, rp)

	return nil
}

func (a *Agent) enablePlugin(name string) error {
	if len(a.PluginNames()) > 0 && !SliceContains(name, a.PluginNames()) {
		return nil
	}

	for _, rp := range a.Plugins {
		if rp.Name == name {
			rp.Enabled = true
		}
	}

	return nil
}

// flusher monitors the events plugin channel and schedules them to correct queues
func (a *Agent) flusher(shutdown chan struct{}, event chan Event) {
	eventCount := 0
	for {
		select {
		case <-shutdown:
			log.Println("Hang on, flushing any cached metrics before shutdown")
			return
		case e := <-event:
			ev_handled := false

			for _, plugin := range a.Plugins {
				subscribedTo := plugin.Plugin.Subscribe(e)
				if SliceContains(e.PayloadModel, subscribedTo) || SliceContains(e.Name, subscribedTo) {
					ev_handled = true
					if a.testRun {
						plugin.Plugin.Process(e)
					} else {
						log.Printf("Enqueue event %v for %v\n", e.Name, plugin.Name)
						workers.Enqueue(plugin.Name, "Event", e)
					}
				}
			}

			if a.testRun && (SliceContains(e.PayloadModel, a.TestSubscribeTo) || SliceContains(e.Name, a.TestSubscribeTo)) {
				ev_handled = true
				eventCount++
				a.TestWork(&e, eventCount)
			}

			if !ev_handled {
				log.Println("Event not handled by any plugin")
				spew.Dump(e)
			}
		}
	}
}

// Run runs the agent daemon
func (a *Agent) Run(shutdown chan struct{}) error {
	var wg sync.WaitGroup

	workers.Middleware = workers.NewMiddleware(
		&workers.MiddlewareRetry{},
		&workers.MiddlewareStats{},
	)

	workers.Configure(map[string]string{
		"server":   viper.GetString("redis.server"),
		"database": viper.GetString("redis.database"),
		"pool":     viper.GetString("redis.pool"),
		"process":  uuid.New(),
	})

	// channel shared between all plugin threads for accumulating events
	a.Events = make(chan Event, 10000)

	for _, plugin := range a.Plugins {
		if !plugin.Enabled {
			continue
		}

		// Start service of any Plugins
		switch p := plugin.Plugin.(type) {
		case Plugin:
			if err := p.Start(a.Events); err != nil {
				log.Printf("Service for plugin %s failed to start, exiting\n%s\n",
					plugin.Name, err.Error())
				return err
			}
			workers.Process(plugin.Name, plugin.Work, viper.GetInt("plugins."+plugin.Name+".workers"))
			defer p.Stop()
		}
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		a.flusher(shutdown, a.Events)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		workers.Run()
		close(shutdown)
	}()

	wg.Wait()
	return nil
}

// Run runs the agent daemon
func (a *Agent) RunTest(shutdown chan struct{}) error {
	var wg sync.WaitGroup
	a.testRun = true

	for _, plugin := range a.Plugins {
		if !plugin.Enabled {
			continue
		}

		// Start service of any Plugins
		switch p := plugin.Plugin.(type) {
		case Plugin:
			if err := p.Start(a.Events); err != nil {
				log.Printf("Service for plugin %s failed to start, exiting\n%s\n",
					plugin.Name, err.Error())
				return err
			}
			defer p.Stop()
		}
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		a.flusher(shutdown, a.Events)
	}()

	wg.Wait()
	return nil
}
