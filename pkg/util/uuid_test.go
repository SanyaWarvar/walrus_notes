package util

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewUUID(t *testing.T) {
	id := NewUUID()
	if CheckUUIDIsZero(id) {
		t.Error("NewUUID() returned zero UUID")
	}
}

func TestCheckUUIDIsZero(t *testing.T) {
	if !CheckUUIDIsZero(uuid.Nil) {
		t.Error("CheckUUIDIsZero(uuid.Nil) = false, want true")
	}
	if CheckUUIDIsZero(NewUUID()) {
		t.Error("CheckUUIDIsZero(non-zero UUID) = true, want false")
	}
}

func TestUUIDFromString(t *testing.T) {
	valid := "550e8400-e29b-41d4-a716-446655440000"

	t.Run("valid UUID", func(t *testing.T) {
		id, err := UUIDFromString(valid)
		if err != nil {
			t.Fatalf("UUIDFromString() error = %v", err)
		}
		if id.String() != valid {
			t.Errorf("UUIDFromString() = %q, want %q", id.String(), valid)
		}
	})

	t.Run("invalid UUID", func(t *testing.T) {
		_, err := UUIDFromString("not-a-uuid")
		if err == nil {
			t.Fatal("expected error for invalid UUID")
		}
	})
}
