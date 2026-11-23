package grpcadapter

import (
	userv1 "contracts/generated/user/v1"
	"user-service/internal/core/dto"
)

func toLogoutInput(req *userv1.LogoutRequest) dto.LogoutInput {
	return dto.LogoutInput{RefreshToken: req.RefreshToken}
}

func toLoginInput(req *userv1.LoginRequest) dto.LoginInput {
	return dto.LoginInput{
		Identifier: req.LoginOrEmail,
		Password:   req.Password,
	}
}

func toLoginResponse(output dto.LoginOutput) *userv1.LoginResponse {
	return &userv1.LoginResponse{
		AccessToken:      output.AccessToken,
		RefreshToken:     output.RefreshToken,
		AccessExpiresIn:  int64(output.AccessExpiresIn.Seconds()),
		RefreshExpiresIn: int64(output.RefreshExpiresIn.Seconds()),
	}
}

func toRegisterInput(req *userv1.RegisterRequest) dto.RegisterInput {
	return dto.RegisterInput{
		Email:    req.Email,
		Login:    req.Login,
		Password: req.Password,
	}
}

func toRegisterResponse(output dto.RegisterOutput) *userv1.RegisterResponse {
	return &userv1.RegisterResponse{
		UserId:           output.UserID,
		AccessToken:      output.Session.AccessToken,
		RefreshToken:     output.Session.RefreshToken,
		AccessExpiresIn:  int64(output.Session.AccessExpiresIn.Seconds()),
		RefreshExpiresIn: int64(output.Session.RefreshExpiresIn.Seconds()),
	}
}

func toRefreshInput(req *userv1.RefreshRequest) dto.RefreshInput {
	return dto.RefreshInput{
		RefreshToken: req.RefreshToken,
	}
}

func toRefreshResponse(output dto.RefreshOutput) *userv1.RefreshResponse {
	return &userv1.RefreshResponse{
		AccessToken:      output.AccessToken,
		RefreshToken:     output.RefreshToken,
		AccessExpiresIn:  int64(output.AccessExpiresIn.Seconds()),
		RefreshExpiresIn: int64(output.RefreshExpiresIn.Seconds()),
	}
}
