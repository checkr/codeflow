package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/codeamp/logger"
	"github.com/codeamp/transistor"

	_ "github.com/codeamp/transistor/example/plugins"
)

func main() {
	config := transistor.Config{
		Server:   "0.0.0.0:16379",
		Database: "0",
		Pool:     "30",
		Process:  "1",
		Queueing: true,
		Plugins: map[string]interface{}{
			"examplePlugin1": map[string]interface{}{
				"hello":   "world1",
				"workers": 1,
			},
			"examplePlugin2": map[string]interface{}{
				"hello":   "world2",
				"workers": 1,
			},
		},
		EnabledPlugins: []string{"examplePlugin1", "examplePlugin2"},
	}

	t, err := transistor.NewTransistor(config)
	if err != nil {
		log.Fatal(err)
	}

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-signals
		if sig == os.Interrupt || sig == syscall.SIGTERM {
			log.Info("Shutting down circuit. SIGTERM recieved!\n")
			// If Queueing is ON then workers are responsible for closing Shutdown chan
			if !t.Config.Queueing {
				t.Stop()
			}
		}
	}()

	log.InfoWithFields("plugins loaded", log.Fields{
		"plugins": strings.Join(t.PluginNames(), ","),
	})

	t.Run()
}
