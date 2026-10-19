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
)

var _ userv1.UserServiceServer = (*Server)(nil)

type Server struct {
	userv1.UnimplementedUserServiceServer
	inbPort inbound.Register
}

func NewServer(uc inbound.Register) *Server {
	return &Server{inbPort: uc}
}

func (s *Server) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	input := dto.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
	}

	output, err := s.inbPort.Register(ctx, input)
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
