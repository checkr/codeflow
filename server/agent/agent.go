package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/jrallison/go-workers"
	"github.com/pborman/uuid"
	"github.com/spf13/viper"
)

// Agent runs codeflow and collects data based on the given config
type Agent struct {
	Queueing   bool
	Events     chan Event
	TestEvents chan Event
	Shutdown   chan struct{}
	Plugins    []*RunningPlugin
}

// NewAgent returns an Agent struct based off the given Config
func NewAgent() (*Agent, error) {
	if len(viper.GetStringMap("plugins")) == 0 {
		log.Fatalf("Error: no plugins found, did you provide a valid config file?")
	}

	agent := &Agent{}

	// channel shared between all plugin threads for accumulating events
	agent.Events = make(chan Event, 10000)

	// channel shared between all plugin threads to trigger shutdown
	agent.Shutdown = make(chan struct{})

	if err := agent.LoadPlugins(); err != nil {
		log.Fatal(err)
	}

	return agent, nil
}

// NewTestAgent returns an Agent struct based off the given Config
func NewTestAgent(config []byte) (*Agent, error) {
	var err error
	var agent *Agent

	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(config))

	if agent, err = NewAgent(); err != nil {
		log.Fatalf("Error while initializing agent: %v", err)
	}

	agent.TestEvents = make(chan Event, 10000)
	agent.Queueing = false

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

		// For debugging purposes
		event.Dump()

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

	if viper.GetInt("plugins."+name+".workers") <= 0 {
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
func (a *Agent) flusher() {
	for {
		select {
		case <-a.Shutdown:
			log.Println("Hang on, flushing any cached metrics before shutdown")
			return
		case e := <-a.Events:
			ev_handled := false

			for _, plugin := range a.Plugins {
				if plugin.Enabled {
					subscribedTo := plugin.Plugin.Subscribe()
					if SliceContains(e.PayloadModel, subscribedTo) || SliceContains(e.Name, subscribedTo) {
						ev_handled = true
						if a.Queueing {
							log.Printf("Enqueue event %v for %v\n", e.Name, plugin.Name)
							workers.Enqueue(plugin.Name, "Event", e)
						} else {
							plugin.Plugin.Process(e)
						}
					}
				}
			}

			if a.TestEvents != nil {
				a.TestEvents <- e
			} else if !ev_handled {
				log.Printf("Event not handled by any plugin: %s\n", e.Name)
				spew.Dump(e)
			}
		}
	}
}

// Run runs the agent daemon
func (a *Agent) Run() error {
	var wg sync.WaitGroup

	if a.Queueing {
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
	}

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

			if a.Queueing {
				workers.Process(plugin.Name, plugin.Work, viper.GetInt("plugins."+plugin.Name+".workers"))
			}
		}
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		a.flusher()
	}()

	if a.Queueing {
		wg.Add(1)
		go func() {
			defer wg.Done()
			workers.Run()
			a.Stop()
		}()

		//wg.Add(1)
		//go func() {
		//	defer wg.Done()
		//	workers.StatsServer(8080)
		//}()
	}

	wg.Wait()
	return nil
}

// Shutdown the agent daemon
func (a *Agent) Stop() {
	close(a.Shutdown)
}

// GetTestEvent listens and returns requested event
func (a *Agent) GetTestEvent(name string, timeout time.Duration) Event {
	// timeout in the case that we don't get requested event
	timer := time.NewTimer(time.Second * timeout)
	go func() {
		<-timer.C
		a.Stop()
		log.Fatalf("Timer expired waiting for event: %v", name)
	}()

	for e := range a.TestEvents {
		if e.Name == name {
			timer.Stop()
			return e
		}
	}

	return Event{}
}
