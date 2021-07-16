package server

import (
	"context"

	"github.com/KevinWu0904/crond/proto/types"
)

// CrondGRPCService serves crond gRPC protocol APIs.
type CrondGRPCService struct {
	types.UnimplementedCrondServer
}

// NewCrondGRPCService creates CrondGRPCService.
func NewCrondGRPCService() *CrondGRPCService {
	return &CrondGRPCService{}
}

// SetJob provides gRPC API for users to create or update a job.
func (s *CrondGRPCService) SetJob(ctx context.Context, req *types.SetJobRequest) (*types.SetJobResponse, error) {
	panic("implement me")
}

// GetJob provides gRPC API for users to search a job.
func (s *CrondGRPCService) GetJob(ctx context.Context, req *types.GetJobRequest) (*types.GetJobResponse, error) {
	panic("implement me")
}

// DeleteJob provides gRPC API for users to delete a job.
func (s *CrondGRPCService) DeleteJob(ctx context.Context, req *types.DeleteJobRequest) (*types.DeleteJobResponse, error) {
	panic("implement me")
}
