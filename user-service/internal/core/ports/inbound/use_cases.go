package inbound

import (
	"context"

	"user-service/internal/core/dto"
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
