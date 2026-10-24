package corerrors

import "errors"

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUserNotFound = errors.New("user not found")

// Сессия (refresh-токен) не найдена: истекла, отозвана или никогда не существовала.
var ErrSessionNotFound = errors.New("session not found")
