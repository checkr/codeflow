package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/checkr/codeflow/server/agent"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "github.com/checkr/codeflow/server/plugins"
	_ "github.com/checkr/codeflow/server/plugins/codeflow"
	_ "github.com/checkr/codeflow/server/plugins/docker_build"
	_ "github.com/checkr/codeflow/server/plugins/heartbeat"
	_ "github.com/checkr/codeflow/server/plugins/kubedeploy"
	_ "github.com/checkr/codeflow/server/plugins/webhooks"
	_ "github.com/checkr/codeflow/server/plugins/websockets"
)

var cfgFile string

func main() {
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

	cmdServer.Flags().StringSliceP("run", "r", []string{}, "run plugins a,b,c")
	viper.BindPFlags(cmdServer.Flags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigType("yaml")
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("CF")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

var cmdServer = &cobra.Command{
	Use:  "server [command]",
	Long: `...`,
	Run: func(cmd *cobra.Command, args []string) {
		ag, err := agent.NewAgent()
		if err != nil {
			log.Fatal(err)
		}

		shutdown := make(chan struct{})
		signals := make(chan os.Signal)
		signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
		go func() {
			sig := <-signals
			if sig == os.Interrupt || sig == syscall.SIGTERM {
				log.Printf("Shutting down Codeflow. SIGTERM recieved!\n")
			}
		}()

		log.Printf("Loaded plugins: %s", strings.Join(ag.PluginNames(), " "))

		ag.Run(shutdown)
	},
}
