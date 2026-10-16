package hasher

import (
	"user-service/internal/core/ports/outbound"

	"golang.org/x/crypto/bcrypt"
)

var _ outbound.PasswordHasher = BcryptHasher{}

type BcryptHasher struct{}

func (bh BcryptHasher) Hash(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}

func (bh BcryptHasher) Verify(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
