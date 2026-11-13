package dto

import "time"

type RefreshInput struct {
	RefreshToken string
}

type RefreshOutput struct {
	AccessToken  string
	RefreshToken string
	// Время жизни выданных токенов.
	AccessExpiresIn  time.Duration
	RefreshExpiresIn time.Duration
}
