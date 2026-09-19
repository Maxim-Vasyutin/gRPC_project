package core

type RegisterInput struct {
	Email    string
	Password string
}

type RegisterOutput struct {
	UserID string
}
