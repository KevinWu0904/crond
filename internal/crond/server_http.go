package crond

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

// HTTPServer serves crond HTTP/1.x protocol APIs.
type HTTPServer struct {
}

// NewHTTPServer creates HTTPServer.
func NewHTTPServer() *HTTPServer {
	return &HTTPServer{}
}

// CreateJob provides HTTP API for users to create a job.
func (hs *HTTPServer) CreateJob(c *gin.Context) {
	panic("implement me")
}

// DeleteJob provides HTTP API for users to delete a job.
func (hs *HTTPServer) DeleteJob(c *gin.Context) {
	panic("implement me")
}

// GetJob provides HTTP API for users to get a job.
func (hs *HTTPServer) GetJob(c *gin.Context) {
	panic("implement me")
}

// UpdateJob provides HTTP API for users to update a job.
func (hs *HTTPServer) UpdateJob(c *gin.Context) {
	panic("implement me")
}

// RegisterCrondHTTPServer defines crond HTTP routers.
// nolint
func RegisterCrondHTTPServer(r *gin.Engine, server *HTTPServer) {
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
