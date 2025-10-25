package handler

import (
	"context"

	"github.com/wutthichod/sa-connext/services/user-service/internal/service"
	pb "github.com/wutthichod/sa-connext/shared/proto/user"
	"google.golang.org/grpc"
)

type gRPCHandler struct {
	pb.UnimplementedUserServiceServer
	service service.Service
}

func NewGRPCHandler(server *grpc.Server, service service.Service) *gRPCHandler {
	handler := &gRPCHandler{
		service: service,
	}
	pb.RegisterUserServiceServer(server, handler)
	return handler
}

func (h *gRPCHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	jwtToken, err := h.service.CreateUser(ctx, req)
	if err != nil {
		return &pb.CreateUserResponse{
			Success: false,
		}, err
	}

	return &pb.CreateUserResponse{
		Success:  true,
		JwtToken: *jwtToken,
	}, nil
}

func (h *gRPCHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	jwtToken, err := h.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		return &pb.LoginResponse{
			Success: false,
		}, err
	}

	return &pb.LoginResponse{
		Success:  true,
		JwtToken: *jwtToken,
	}, nil
}
