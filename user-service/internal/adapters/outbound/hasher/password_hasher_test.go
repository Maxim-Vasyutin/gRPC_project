package hasher

import (
	"testing"
)

func TestBcryptHasher_HashAndVerify(t *testing.T) {
	hasher := BcryptHasher{}
	password := "supersecret123"

	hash, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("error hashing: %v", err)
	}

	err = hasher.Verify(password, hash)
	if err != nil {
		t.Fatalf("error verify: %v", err)
	}

	err = hasher.Verify("wrong-password", hash)
	if err == nil {
		t.Fatal("expected verification error for wrong password")
	}
}
