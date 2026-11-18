package domain

import "strings"

type User struct {
	ID           string
	Email        string
	Login        string
	PasswordHash string
}

func NormalizeEmail(s string) string {
	return strings.ToLower(s)
}

func NormalizeLogin(s string) string {
	return strings.ToLower(s)
}