/*
JWT-токен состоинт из трёх частей через точку xxx.yyy.zzz
первая часть - header - алгортим подписи (HS256)
вторая часть - payload - мои данные
третья часть - signature - результат высчисления с секретом. Секрет всегда на сервере

то есть, сервер собирает 1 и 2, потом сам считает подпись
своим ключом и потом сравнивает
*/

package token

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"user-service/internal/core/dto"
	"user-service/internal/core/ports/outbound"

	"github.com/golang-jwt/jwt/v5"
)

// Проверка при компиляции: JWTIssuer подходит к порту TokenIssuer.
var _ outbound.TokenIssuer = (*JWTIssuer)(nil)

type JWTIssuer struct {
	secret    []byte
	accessTTL time.Duration
}

func NewJWTIssuer(secret []byte, accessTTL time.Duration) *JWTIssuer {
	return &JWTIssuer{
		secret:    secret,
		accessTTL: accessTTL,
	}
}

func (jt *JWTIssuer) Issue(userID string) (dto.TokenPair, error) {
	now := time.Now()
	//Регистрируем токен
	//1-е поле: Кому отдаём
	//2-е: когда выдан
	//3-е: когда истечёт
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(jt.accessTTL)),
	}

	//NewWithClaims - собирает header и payload
	//SignedString - считает подпись и склеивает всё в строку
	access, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jt.secret)
	if err != nil {
		return dto.TokenPair{}, fmt.Errorf("sign access token: %w", err)
	}

	refresh, err := newRefreshToken()
	if err != nil {
		return dto.TokenPair{}, err
	}

	return dto.TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

// newRefreshToken: 32 случайных байта из crypto/rand, превращённые в строку.
func newRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
