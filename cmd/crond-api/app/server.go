package app

import (
	"github.com/KevinWu0904/crond/internal/app/crond-api/handler"
	"github.com/KevinWu0904/crond/pkg/logs"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func NewAPIServerCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "crond-api",
		Long: "The CronD API Server provides REST service for jobs.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return StartAPIServer()
		},
	}
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
