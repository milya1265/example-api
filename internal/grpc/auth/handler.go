package auth

import (
	"context"
	"errors"
	"example1/internal/service"
	"example1/pkg/logger"
	appv1 "example1/protos/gen/go/app"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	appv1.UnimplementedAuthServer
	Service service.Auth
	Logger  logger.Logger
}

func RegisterServerAPI(s *grpc.Server, auth *service.Auth) {
	appv1.RegisterAuthServer(s, &serverAPI{Service: *auth, Logger: logger.Get()})
}

func (s *serverAPI) Login(ctx context.Context, req *appv1.LoginReq) (*appv1.LoginRes, error) {
	if req.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid login")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid password")
	}

	token, err := s.Service.Login(ctx, req.GetLogin(), req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &appv1.LoginRes{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *appv1.RegisterReq) (*appv1.RegisterRes, error) {
	if req.GetLogin() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid login")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid password")
	}

	if req.GetRole() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid role")

	}

	userID, err := s.Service.RegisterNewUser(ctx, req.GetLogin(), req.GetPassword(), req.GetRole())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &appv1.RegisterRes{UserId: userID}, nil

}
func (s *serverAPI) GetRole(ctx context.Context, req *appv1.GetRoleReq) (*appv1.GetRoleRes, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	role, err := s.Service.GetRole(ctx, req.GetUserId())
	if err != nil {
		s.Logger.Error()
		if errors.Is(err, service.UserNotFoundErr) {
			return nil, status.Error(codes.Unauthenticated, service.WrongLoginOrPasswordErr.Error())
		}
		if errors.Is(err, service.WrongLoginOrPasswordErr) {
			return nil, status.Error(codes.Unauthenticated, service.WrongLoginOrPasswordErr.Error())
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &appv1.GetRoleRes{Role: role}, nil
}
