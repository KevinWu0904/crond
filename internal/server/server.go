package server

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/KevinWu0904/crond/pkg/logs"
	"github.com/KevinWu0904/crond/proto/types"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

// Server represents crond server.
type Server struct {
	c            *Config
	grpcServer   *grpc.Server
	httpServer   *http.Server
	raftLayer    *RaftLayer
	mux          cmux.CMux
	grpcListener net.Listener
	httpListener net.Listener
	raftListener net.Listener
}

// NewServer creates crond Server.
func NewServer(c *Config) (*Server, error) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(c.ServerPort))
	if err != nil {
		return nil, err
	}

	// Serve multiple protocols on the same listener.
	mux := cmux.New(listener)

	grpcListener := mux.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpListener := mux.Match(cmux.HTTP1Fast())
	raftListener := mux.Match(cmux.Any())

	// New crond gRPC server.
	grpcServer := grpc.NewServer()
	grpcService := NewCrondGRPCService()
	types.RegisterCrondServer(grpcServer, grpcService)

	// New crond HTTP server.
	router := gin.Default()
	gin.DefaultWriter = logs.GetGinWriter()
	gin.DefaultErrorWriter = logs.GetGinErrorWriter()
	router.Use(ginzap.Ginzap(logs.GetLogger(), "2006-01-02T15:04:05.000Z0700", false))
	router.Use(ginzap.RecoveryWithZap(logs.GetLogger(), true))

	httpService := NewCrondHTTPService()
	RegisterCrondHTTPServer(router, httpService)
	httpServer := &http.Server{Handler: router}

	// New crond raft layer.
	raftLayer := NewRaftLayer(c, raftListener)

	return &Server{
		c:            c,
		grpcServer:   grpcServer,
		httpServer:   httpServer,
		raftLayer:    raftLayer,
		mux:          mux,
		grpcListener: grpcListener,
		httpListener: httpListener,
		raftListener: raftListener,
	}, nil
}

// Run launches crond server.
func (s *Server) Run(ctx context.Context) {
	go s.grpcServer.Serve(s.grpcListener)
	go s.httpServer.Serve(s.httpListener)
	go s.raftLayer.Run()

	logs.CtxInfo(ctx, "CronD starting...: port=%d", s.c.ServerPort)
	s.mux.Serve()
}

// GracefulShutdown stops crond server gracefully.
func (s *Server) GracefulShutdown(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	s.grpcServer.GracefulStop()
	s.httpServer.Shutdown(ctx)

	logs.CtxInfo(ctx, "CronD shutdown gracefully")
}
