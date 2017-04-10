package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/checkr/codeflow/server/agent"
	"github.com/mattes/migrate/file"
	"github.com/mattes/migrate/migrate"
	"github.com/mattes/migrate/migrate/direction"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "github.com/checkr/codeflow/server/plugins"
	_ "github.com/checkr/codeflow/server/plugins/codeflow"
	_ "github.com/checkr/codeflow/server/plugins/codeflow/migrations"
	_ "github.com/checkr/codeflow/server/plugins/docker_build"
	_ "github.com/checkr/codeflow/server/plugins/git_sync"
	_ "github.com/checkr/codeflow/server/plugins/heartbeat"
	_ "github.com/checkr/codeflow/server/plugins/kubedeploy"
	_ "github.com/checkr/codeflow/server/plugins/slack"
	_ "github.com/checkr/codeflow/server/plugins/webhooks"
	_ "github.com/checkr/codeflow/server/plugins/websockets"
)

var cfgFile string

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	Execute()
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use: "codeflow",
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./configs/.codeflow.yml)")
	RootCmd.AddCommand(cmdServer)
	RootCmd.AddCommand(cmdMigrate)
	cmdMigrate.AddCommand(cmdMigrateUp)
	cmdMigrate.AddCommand(cmdMigrateDown)

	cmdServer.Flags().StringSliceP("run", "r", []string{}, "run plugins a,b,c")
	viper.BindPFlags(cmdServer.Flags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("CF")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	viper.AutomaticEnv() // read in environment variables that match
}

var cmdServer = &cobra.Command{
	Use:  "server [command]",
	Long: `...`,
	Run: func(cmd *cobra.Command, args []string) {
		ag, err := agent.NewAgent()
		if err != nil {
			log.Fatal(err)
		}

		ag.Queueing = true

		signals := make(chan os.Signal)
		signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
		go func() {
			sig := <-signals
			if sig == os.Interrupt || sig == syscall.SIGTERM {
				log.Printf("Shutting down Codeflow. SIGTERM recieved!\n")
				// If Queueing is ON then workers are responsible for closing Shutdown chan
				if !ag.Queueing {
					ag.Stop()
				}
			}
		}()

		log.Printf("Loaded plugins: %s", strings.Join(ag.PluginNames(), " "))

		ag.Run()
	},
}

var cmdMigrate = &cobra.Command{
	Use:  "migrate [command]",
	Long: `...`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}
var cmdMigrateUp = &cobra.Command{
	Use:  "up [command]",
	Long: `...`,
	Run: func(cmd *cobra.Command, args []string) {
		pipe := migrate.NewPipe()
		go migrate.Up(pipe, viper.GetString("plugins.codeflow.mongodb.uri"), "./plugins/codeflow/migrations")
		ok := writePipe(pipe)
		if !ok {
			os.Exit(1)
		}
	},
}

var cmdMigrateDown = &cobra.Command{
	Use:  "down [command]",
	Long: `...`,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetString("environment") != "development" {
			panic("You can only use migrate down in development environment")
		}

		pipe := migrate.NewPipe()
		go migrate.Down(pipe, viper.GetString("plugins.codeflow.mongodb.uri"), "./plugins/codeflow/migrations")
		ok := writePipe(pipe)
		if !ok {
			os.Exit(1)
		}
	},
}

func writePipe(pipe chan interface{}) (ok bool) {
	okFlag := true
	if pipe != nil {
		for {
			select {
			case item, more := <-pipe:
				if !more {
					return okFlag
				}
				switch item.(type) {

				case string:
					fmt.Println(item.(string))

				case error:
					fmt.Printf("%s\n\n", item.(error).Error())
					okFlag = false

				case file.File:
					f := item.(file.File)
					if f.Direction == direction.Up {
						fmt.Print(">")
					} else if f.Direction == direction.Down {
						fmt.Print("<")
					}
					fmt.Printf(" %s\n", f.FileName)

				default:
					text := fmt.Sprint(item)
					fmt.Println(text)
				}
			}
		}
	}
	return okFlag
}
