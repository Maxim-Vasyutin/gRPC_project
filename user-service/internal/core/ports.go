package core

import (
	"context"
)

type RegisterUserRepository interface {
	//Проверка наличия email
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	//Сохранения пользователя в бд
	Save(ctx context.Context, user User) error
}

type LoginUserRepository interface {
	//Вытащить пользователя по почте
	GetByEmail(ctx context.Context, email string) (User, error)
}

type PasswordService interface {
	//Хэширование пароля
	Hash(password string) (string, error)
	//Верификация пароля и хэша
	Verify(password, hash string) error
}

type IDGenerator interface {
	//Генерация ID нового пользователя
	NewID() string
}

type TokenIssuer interface {
	//Генерация пары токенов
	Issue(userID string) (TokenPair, error)
}

type SessionRepository interface {
	//Сохранение сессию
	Store(ctx context.Context, session Session) error
	//Перезапись старой сессии
	Replace(ctx context.Context, oldRefreshToken string, newSession Session) error
	//Поиск сессии по рефреш-токену
	FindByRefreshToken(ctx context.Context, refreshToken string) (Session, error)
	//Отзыв сессии
	RevokeByRefreshToken(ctx context.Context, refreshToken string) error
}
