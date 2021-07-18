package cmd

import (
	"fmt"
	"os"

	"github.com/KevinWu0904/crond/internal/server"
	"github.com/KevinWu0904/crond/pkg/logs"
	"github.com/spf13/pflag"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string
var config = DefaultConfig()

// RootCommand represents crond CLI.
var RootCommand = &cobra.Command{
	Use:   "crond",
	Short: "CronD is a Cloud Native golang distributed cron scheduling service",
	Long: `CronD serves a distributed unified job dispatcher for offline periodic tasks. It is recommended running in 
a cluster with 3 or 5 nodes, peer nodes communicates by Raft Consensus`,
}

// init retrieves crond configs from file/env/flag, priority is file > env > flag.
func init() {
	cobra.OnInitialize(initConfig)

	// Add crond sub commands.
	RootCommand.AddCommand(ServerCommand)

	// Bind crond global config file.
	RootCommand.PersistentFlags().StringVarP(&configFile, "config", "c", "", "server global config file")

	// Bind crond extra flags to related commands.
	bindRootFlags()
	bindServerFlags()
}

func bindRootFlags() {
	fs := pflag.NewFlagSet("", pflag.ExitOnError)
	logs.BindFlags(config.Logger, fs)
	RootCommand.PersistentFlags().AddFlagSet(fs)
}

func bindServerFlags() {
	fs := pflag.NewFlagSet("", pflag.ExitOnError)
	server.BindFlags(config.Server, fs)
	ServerCommand.Flags().AddFlagSet(fs)
}

// initConfig reads configs from specific directories or environment variables.
func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("crond-config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.crond")
		viper.AddConfigPath("/etc/crond")
	}

	viper.SetEnvPrefix("crond")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "initConfig failed to load config file: err=%v\n", err)
		return
	}
	fmt.Fprintf(os.Stdout, "initConfig succeed to load config file: file=%s\n", viper.ConfigFileUsed())

	if err := viper.Unmarshal(config); err != nil {
		fmt.Fprintf(os.Stderr, "initConfig failed to unmarshal config file: err=%v", err)
	}
}
