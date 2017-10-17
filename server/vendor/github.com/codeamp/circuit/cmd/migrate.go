package cmd

import (
	log "github.com/codeamp/logger"
	"github.com/codeamp/transistor"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run plugin migrations",
	Run: func(cmd *cobra.Command, args []string) {
		config := transistor.Config{
			Queueing:       false,
			Plugins:        viper.GetStringMap("plugins"),
			EnabledPlugins: viper.GetStringSlice("enable"),
		}

		t, err := transistor.NewTransistor(config)
		if err != nil {
			log.Fatal(err)
		}

		for _, plugin := range t.Plugins {
			if !plugin.Enabled {
				continue
			}

			switch p := plugin.Plugin.(type) {
			case transistor.Plugin:
				if _p, ok := p.(interface {
					Migrate()
				}); ok {
					_p.Migrate()
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(migrateCmd)
}
