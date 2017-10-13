package cmd

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	log "github.com/codeamp/logger"
	"github.com/codeamp/transistor"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startCmd = &cobra.Command{
	Use:  "start",
	Long: `...`,
	Run: func(cmd *cobra.Command, args []string) {
		config := transistor.Config{
			Server:         viper.GetString("redis.server"),
			Database:       viper.GetString("redis.database"),
			Pool:           viper.GetString("redis.pool"),
			Process:        viper.GetString("redis.process"),
			Queueing:       true,
			Plugins:        viper.GetStringMap("plugins"),
			EnabledPlugins: viper.GetStringSlice("enable"),
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
				log.Info("Shutting down Circuit. SIGTERM recieved!\n")
				// If Queueing is ON then workers are responsible for closing Shutdown chan
				if !config.Queueing {
					t.Stop()
				}
			}
		}()

		log.InfoWithFields("plugins loaded", log.Fields{
			"plugins": strings.Join(t.PluginNames(), ","),
		})

		t.Run()
	},
}

func init() {
	RootCmd.AddCommand(startCmd)
	startCmd.Flags().StringSliceP("enable", "e", []string{}, "websockets")
	viper.BindPFlags(startCmd.Flags())
}
