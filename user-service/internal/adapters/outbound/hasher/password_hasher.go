package hasher

import (
	"fmt"

	"user-service/internal/core/ports/outbound"

	"golang.org/x/crypto/bcrypt"
)

var _ outbound.PasswordHasher = BcryptHasher{}

// BcryptHasher хеширует пароли bcrypt-ом с заданной стоимостью (cost).
// Создавать через NewBcryptHasher: конструктор проверяет cost.
type BcryptHasher struct {
	cost int
}

// NewBcryptHasher проверяет, что cost попадает в допустимый диапазон bcrypt.
func NewBcryptHasher(cost int) (BcryptHasher, error) {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return BcryptHasher{}, fmt.Errorf("bcrypt cost must be in [%d, %d], got %d", bcrypt.MinCost, bcrypt.MaxCost, cost)
	}
	return BcryptHasher{cost: cost}, nil
}

func (bh BcryptHasher) Hash(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bh.cost)
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
