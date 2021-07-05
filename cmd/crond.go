package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/KevinWu0904/crond/internal/crond"
	"github.com/KevinWu0904/crond/pkg/flag"
	"github.com/KevinWu0904/crond/pkg/logs"
	"github.com/KevinWu0904/crond/pkg/term"
	"github.com/KevinWu0904/crond/proto/types"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/soheilhy/cmux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var config = crond.DefaultConfig()
var configFile string

// Command represents the crond CLI.
var Command = &cobra.Command{
	Use:   "crond",
	Short: "CronD is a Cloud Native golang distributed cron scheduling service.",
	Long: `CronD serves a distributed unified job dispatcher for offline periodic tasks. It is recommended running in 
a cluster with 3 or 5 nodes, peer nodes communicates by Raft Consensus.`,
	Run: Run,
}

func init() {
	cobra.OnInitialize(initConfig)

	Command.PersistentFlags().StringVar(&configFile, "config", "", "crond config")

	nfs := flag.NewNamedFlagSets()

	// Bind custom named flag sets.
	logs.BindLoggerFlags(config.Logger, nfs.NewFlatSet("logger"))
	crond.BindServerFlags(config.Server, nfs.NewFlatSet("server"))

	for _, fs := range nfs.FlagSets {
		Command.Flags().AddFlagSet(fs)
	}
	viper.BindPFlags(Command.Flags())

	// Custom crond CLI usage and help.
	usageTpl := "Usage:\n  %s\n"
	cols, _, _ := term.TerminalSize(Command.OutOrStdout())
	Command.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageTpl, cmd.UseLine())
		flag.PrintSections(cmd.OutOrStderr(), nfs, cols)
		return nil
	})

	helpTpl := "Name:\n  %s\n\nDescription:\n  %s\n\n" + usageTpl
	Command.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), helpTpl, cmd.Short, cmd.Long, cmd.UseLine())
		flag.PrintSections(cmd.OutOrStdout(), nfs, cols)
	})
}

// initConfig reads configs from specific directories or environment variables.
func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName("crond-config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("/etc/crond")
		viper.AddConfigPath("$HOME/.crond")
		viper.AddConfigPath("./conf")
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

// Run starts crond servers.
func Run(cmd *cobra.Command, args []string) {
	if err := logs.InitLogger(config.Logger); err != nil {
		panic(err)
	}
	defer logs.Flush()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer stop()

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(config.Server.ServerPort))
	if err != nil {
		panic(err)
	}

	logs.CtxInfo(ctx, "CronD start server ...: port=%d", config.Server.ServerPort)

	mux := cmux.New(listener)

	grpcListener := mux.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpListener := mux.Match(cmux.HTTP1Fast())

	grpcServer := grpc.NewServer()
	crondGRPCServer := crond.NewGRPCServer()
	types.RegisterCrondServer(grpcServer, crondGRPCServer)
	logs.CtxInfo(ctx, "CronD start gRPC Server ...")
	go grpcServer.Serve(grpcListener)

	router := gin.Default()
	router.Use(ginzap.Ginzap(logs.GetLogger(), "2006-01-02T15:04:05.000Z0700", false))
	router.Use(ginzap.RecoveryWithZap(logs.GetLogger(), true))

	httpServer := &http.Server{Handler: router}
	crondHTTPServer := crond.NewHTTPServer()
	crond.RegisterCrondHTTPServer(router, crondHTTPServer)
	logs.CtxInfo(ctx, "CronD start HTTP Server ...")
	go httpServer.Serve(httpListener)

	logs.CtxInfo(ctx, "CronD servers starting ...")
	go mux.Serve()

	<-ctx.Done()
	stop()

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	grpcServer.GracefulStop()
	httpServer.Shutdown(ctx)

	logs.CtxInfo(ctx, "CronD stop gracefully")
}
