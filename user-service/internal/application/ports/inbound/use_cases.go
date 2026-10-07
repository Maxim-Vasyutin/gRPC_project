package inbound

import (
	"context"

	"gRPC_project/user-service/internal/application/dto"
)

type Register interface {
	Register(context.Context, dto.RegisterInput) (dto.RegisterOutput, error)
}

type Login interface {
	Login(context.Context, dto.LoginInput) (dto.LoginOutput, error)
}

type Refresh interface {
	Refresh(context.Context, dto.RefreshInput) (dto.RefreshOutput, error)
}

type Logout interface {
	Logout(context.Context, dto.LogoutInput) error
}
