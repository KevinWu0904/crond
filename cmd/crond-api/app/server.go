package app

import (
	"fmt"

	"github.com/KevinWu0904/crond/internal/app/crond-api/handler"
	"github.com/KevinWu0904/crond/pkg/flag"
	"github.com/KevinWu0904/crond/pkg/logs"
	"github.com/KevinWu0904/crond/pkg/term"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

type APIServerOptions struct {
	LoggerOptions *logs.LoggerOptions
}

func NewAPIServerOptions() *APIServerOptions {
	return &APIServerOptions{
		LoggerOptions: logs.NewLoggerOptions(),
	}
}

func (options *APIServerOptions) BindNamedFlagSets(nfs *flag.NamedFlagSets) {
	options.LoggerOptions.BindFlags(nfs.NewFlatSet("log"))
}

func NewAPIServerCommand(options *APIServerOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "crond-api",
		Long: "The CronD API Server provides REST service for jobs.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return StartAPIServer()
		},
	}

	fs := cmd.Flags()
	nfs := flag.NewNamedFlagSets()
	options.BindNamedFlagSets(nfs)
	for _, f := range nfs.FlagSets {
		fs.AddFlagSet(f)
	}

	usageFmt := "Usage:\n  %s\n"
	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		flag.PrintSections(cmd.OutOrStderr(), nfs, cols)
		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		flag.PrintSections(cmd.OutOrStdout(), nfs, cols)
	})

	return cmd
}

func StartAPIServer() error {
	r := gin.New()

	v1 := r.Group("/v1")
	{
		v1.GET("/jobs", handler.GetJobs)
		v1.POST("/jobs", handler.CreateJob)
		v1.PUT("/jobs/:job_id", handler.UpdateJob)
		v1.DELETE("/jobs/:job_id", handler.DeleteJob)
	}

	logs.Info("StartAPIServer ...")

	if err := r.Run(); err != nil {
		logs.Error("StartAPIServer failed: err=%v", err)
		return err
	}

	return nil
}
