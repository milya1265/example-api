package grpcserver

import (
	authgrc "example1/internal/grpc/auth"
	"example1/internal/service"
	"example1/pkg/logger"
	"google.golang.org/grpc"
	"net"
	"time"
)

type gRPCSrv struct {
	GrpcSrv *grpc.Server
	Logger  logger.Logger
}

type GRPCServer interface {
	Serve(l net.Listener, auth *service.Auth) error
}

func (g *gRPCSrv) Serve(l net.Listener, auth *service.Auth) error {

	authgrc.RegisterServerAPI(g.GrpcSrv, auth)

	err := g.GrpcSrv.Serve(l)

	return err
}

func New(
	grpcPort string,
	tokenTTL time.Duration,
) GRPCServer {

	s := &gRPCSrv{
		GrpcSrv: grpc.NewServer(),
		Logger:  logger.Get(),
	}

	return s
}
