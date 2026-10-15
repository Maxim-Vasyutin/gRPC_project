package postgres

import (
	"context"
	"testing"
	"user-service/internal/core"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresUserRepository_SaveAndGetByEmail(t *testing.T) {
	ctx := context.Background()
	dsn := "postgres://user:password@localhost:5434/users?sslmode=disable"

	//создаём соединение
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	//проверяем соединение
	if err := pool.Ping(ctx); err != nil {
		t.Fatal(err)
	}

	//создаём экземпляр структуры, чтобы с ним работать дальше
	repo := NewPostgresUserRepository(pool)

	//создаём тестового пользователя
	user := core.User{
		ID:           "550e8400-e29b-41d4-a716-446655440000",
		Email:        "test@example.com",
		PasswordHash: "some-hash",
	}

	//если пользователь найден - удаляем его
	_, err = pool.Exec(
		ctx,
		`DELETE FROM users WHERE email = $1`,
		user.Email,
	)
	if err != nil {
		t.Fatal(err)
	}

	//фиксируем defer для очистки данных из БД
	defer func() {
		_, _ = pool.Exec(ctx, `DELETE FROM users WHERE email = $1`, user.Email)
	}()

	//сохраняем пользователя, которого создали выше
	if err = repo.Save(ctx, user); err != nil {
		t.Fatal(err)
	}

	//Проверяем, получилось ли создать пользователя. Запрашиваем его назад
	got, err := repo.GetByEmail(ctx, user.Email)
	if err != nil {
		t.Fatal(err)
	}

	//Сравниваем 2 структуры пользователей
	if got != user {
		t.Fatalf("got %+v, want %+v", got, user)
	}

}
