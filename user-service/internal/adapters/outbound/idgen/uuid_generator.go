package idgen

import (
	"user-service/internal/application/ports/outbound"

	"github.com/google/uuid"
)

var _ outbound.IDGenerator = UUIDGenerator{}

type UUIDGenerator struct{}

func (g UUIDGenerator) NewID() string {
	return uuid.NewString()
}
