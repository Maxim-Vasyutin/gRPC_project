package dto

import "time"

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	AccessToken  string
	RefreshToken string
	// Время жизни выданных токенов.
	AccessExpiresIn  time.Duration
	RefreshExpiresIn time.Duration
}
