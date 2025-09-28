package util

import (
	"github.com/google/uuid"
)

func NewUUID() uuid.UUID {
	return uuid.New()
}

func CheckUUIDIsZero(id uuid.UUID) bool {
	return id == uuid.Nil
}

func UUIDFromString(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}
