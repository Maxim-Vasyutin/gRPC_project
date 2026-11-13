package dto

import "time"

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	// AccessTTL — время жизни access-токена, с которым он реально выпущен.
	AccessTTL time.Duration
}
