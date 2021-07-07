package crond

import (
	"context"

	"github.com/KevinWu0904/crond/proto/types"
)

// GRPCServer serves crond gRPC protocol APIs.
type GRPCServer struct {
	types.UnimplementedCrondServer
}

// NewGRPCServer creates GRPCServer.
func NewGRPCServer() *GRPCServer {
	return &GRPCServer{}
}

// SetJob provides gRPC API for users to create or update a job.
func (s *GRPCServer) SetJob(ctx context.Context, req *types.SetJobRequest) (*types.SetJobResponse, error) {
	panic("implement me")
}

// GetJob provides gRPC API for users to search a job.
func (s *GRPCServer) GetJob(ctx context.Context, req *types.GetJobRequest) (*types.GetJobResponse, error) {
	panic("implement me")
}

// DeleteJob provides gRPC API for users to delete a job.
func (s *GRPCServer) DeleteJob(ctx context.Context, req *types.DeleteJobRequest) (*types.DeleteJobResponse, error) {
	panic("implement me")
}
