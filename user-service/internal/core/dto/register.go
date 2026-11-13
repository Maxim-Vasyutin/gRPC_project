package dto

type RegisterInput struct {
	Email    string
	Login    string
	Password string
}

type RegisterOutput struct {
	UserID string
}
