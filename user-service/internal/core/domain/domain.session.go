package domain

import "time"

type Session struct {
	UserID       string
	RefreshToken string
	ExpiresAt    time.Time
}
