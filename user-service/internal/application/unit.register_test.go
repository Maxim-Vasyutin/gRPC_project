package application

import (
	"context"
	"errors"
	"testing"

	"gRPC_project/user-service/internal/application/dto"
	"gRPC_project/user-service/internal/core"
)

type fakeUserRepository struct {
	existsErr  error
	saveCalled bool
}

type fakePasswordHasher struct {
	called bool
}

type fakeIDGenerator struct {
	called bool
}

func TestRegister_RepositoryCheckError(t *testing.T) {
	//average
	repoErr := errors.New("repository error")

	repo := &fakeUserRepository{
		existsErr: repoErr,
	}
	hasher := &fakePasswordHasher{}
	idGenerator := &fakeIDGenerator{}

	//act -
	useCase := NewRegisterUseCase(repo, hasher, idGenerator)

	ctx := context.Background()

	//assert
	_, err := useCase.Register(ctx, dto.RegisterInput{
		Email:    "test@mail.com",
		Password: "12344321",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if hasher.called {
		t.Error("hasher should not be called")
	}

	if idGenerator.called {
		t.Error("idGenerator should not be called")
	}

	if repo.saveCalled {
		t.Error("save should not be called")
	}

}

// fakeUserRepository
func (ur *fakeUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {

	return false, ur.existsErr
}

func (ur *fakeUserRepository) Save(ctx context.Context, user core.User) error {

	ur.saveCalled = true
	return nil
}

func (ur *fakeUserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (core.User, error) {
	return core.User{}, nil
}

// fakePasswordHasher
func (h *fakePasswordHasher) Hash(password string) (string, error) {
	h.called = true
	return "hash", nil
}
func (h *fakePasswordHasher) Verify(password, hash string) error {
	return nil
}

// fakeIDGenerator
func (g *fakeIDGenerator) NewID() string {
	g.called = true
	return "test-id"
}
