package cmd

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/KevinWu0904/crond/internal/server"
	"github.com/KevinWu0904/crond/pkg/logs"
	"github.com/spf13/cobra"
)

// ServerCommand represents crond server CLI.
var ServerCommand = &cobra.Command{
	Use:   "server",
	Short: "CronD server is the actual real server endpoint",
	Long:  "CronD server actually launches the HA distributed cron scheduling servers",
	Run:   RunServer,
}

// RunServer launches crond server.
func RunServer(cmd *cobra.Command, args []string) {
	if err := logs.InitLogger(config.RootConfig.Logger); err != nil {
		panic(err)
	}
	defer logs.Flush()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer stop()

	cs, err := server.NewServer(config.ServerConfig.Server)
	if err != nil {
		panic(err)
	}

	go cs.Run(ctx)

	<-ctx.Done()

	cs.GracefulShutdown(ctx)
}
