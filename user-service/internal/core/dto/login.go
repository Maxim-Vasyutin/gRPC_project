package dto

import "time"

type LoginInput struct {
	// Identifier — email или логин (если содержит "@", это email).
	Identifier string
	Password   string
}

type LoginOutput struct {
	AccessToken  string
	RefreshToken string
	// Время жизни выданных токенов.
	AccessExpiresIn  time.Duration
	RefreshExpiresIn time.Duration
}
