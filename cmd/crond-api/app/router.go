package app

import (
	"github.com/KevinWu0904/crond/internal/app/crond-api/handler"
	"github.com/KevinWu0904/crond/pkg/logs"
	"github.com/gin-gonic/gin"
)

func RegisterAPIServer(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		v1.GET("/jobs", handler.GetJobs)
		v1.POST("/jobs", handler.CreateJob)
		v1.PUT("/jobs/:job_id", handler.UpdateJob)
		v1.DELETE("/jobs/:job_id", handler.DeleteJob)
	}

	logs.Info("RegisterAPIServer success")
}
