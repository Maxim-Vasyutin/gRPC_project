//Принимает запрос из protobuf request
//Вытаскивает данные. Отправляет в core. И формирует ответ.

package grpcadapter

import (
	"context"
	userv1 "contracts/generated/user/v1"
	"errors"
	"log"
	"user-service/internal/core/corerrors"
	"user-service/internal/core/dto"
	"user-service/internal/core/ports/inbound"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ userv1.UserServiceServer = (*Server)(nil)

type Server struct {
	userv1.UnimplementedUserServiceServer
	register inbound.Register
	login    inbound.Login
	refresh  inbound.Refresh
	logout   inbound.Logout
}

func NewServer(
	register inbound.Register,
	login inbound.Login,
	refresh inbound.Refresh,
	logout inbound.Logout,
) *Server {
	return &Server{
		register: register,
		login:    login,
		refresh:  refresh,
		logout:   logout,
	}
}

func (s *Server) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	input := dto.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := s.register.Register(ctx, input)
	if err != nil {
		if errors.Is(err, corerrors.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		log.Printf("register: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &userv1.RegisterResponse{
		UserId: output.UserID,
	}, nil
}

func (s *Server) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	input := dto.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := s.login.Login(ctx, input)
	if err != nil {
		if errors.Is(err, corerrors.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid email or password")
		}
		log.Printf("login: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &userv1.LoginResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}, nil
}

func (s *Server) Refresh(ctx context.Context, req *userv1.RefreshRequest) (*userv1.RefreshResponse, error) {
	input := dto.RefreshInput{
		RefreshToken: req.RefreshToken,
	}

	output, err := s.refresh.Refresh(ctx, input)
	if err != nil {
		if errors.Is(err, corerrors.ErrSessionNotFound) {
			return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
		}
		log.Printf("refresh: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &userv1.RefreshResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}, nil
}

func (s *Server) Logout(ctx context.Context, req *userv1.LogoutRequest) (*emptypb.Empty, error) {
	input := dto.LogoutInput{
		RefreshToken: req.RefreshToken,
	}

	err := s.logout.Logout(ctx, input)
	if err != nil {
		log.Printf("logout: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &emptypb.Empty{}, nil
}
