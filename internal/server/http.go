package server

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

// CrondHTTPService serves crond HTTP/1.x protocol APIs.
type CrondHTTPService struct {
}

// NewCrondHTTPService creates CrondHTTPService.
func NewCrondHTTPService() *CrondHTTPService {
	return &CrondHTTPService{}
}

// CreateJob provides HTTP API for users to create a job.
func (hs *CrondHTTPService) CreateJob(c *gin.Context) {
	panic("implement me")
}

// DeleteJob provides HTTP API for users to delete a job.
func (hs *CrondHTTPService) DeleteJob(c *gin.Context) {
	panic("implement me")
}

// GetJob provides HTTP API for users to get a job.
func (hs *CrondHTTPService) GetJob(c *gin.Context) {
	panic("implement me")
}

// UpdateJob provides HTTP API for users to update a job.
func (hs *CrondHTTPService) UpdateJob(c *gin.Context) {
	panic("implement me")
}

// RegisterCrondHTTPServer registers crond HTTP routers.
// nolint
func RegisterCrondHTTPServer(r *gin.Engine, server *CrondHTTPService) {
	pprof.Register(r)

	v1 := r.Group("/v1")
	{
		jobs := v1.Group("/jobs")
		{
			jobs.POST("", server.CreateJob)
			jobs.DELETE("/:job_id", server.DeleteJob)
			jobs.GET("/:job_id", server.GetJob)
			jobs.PUT("/:job_id", server.UpdateJob)
		}
	}
}
