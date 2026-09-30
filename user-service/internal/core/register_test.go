package core

import (
	"context"
	"errors"
	"testing"
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
	_, err := useCase.Register(ctx, RegisterInput{
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

func (ur *fakeUserRepository) Save(ctx context.Context, user User) error {

	ur.saveCalled = true
	return nil
}

// fakePasswordHasher
func (h *fakePasswordHasher) Hash(password string) (string, error) {
	h.called = true
	return "hash", nil
}

// fakeIDGenerator
func (g *fakeIDGenerator) NewID() string {
	g.called = true
	return "test-id"
}
