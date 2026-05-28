package entity

import (
	"testing"
	"wn/internal/domain/services/crypto"

	"github.com/google/uuid"
)

func TestNote_GetId(t *testing.T) {
	id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	n := Note{Id: id}
	if n.GetId() != id {
		t.Errorf("GetId() = %v, want %v", n.GetId(), id)
	}
}

func TestNoteWithPosition_GetId(t *testing.T) {
	id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	n := NoteWithPosition{Note: Note{Id: id}}
	if n.GetId() != id {
		t.Errorf("GetId() = %v, want %v", n.GetId(), id)
	}
}

func TestNote_EncryptDecrypt(t *testing.T) {
	encryptor := crypto.NewEncryptor("test-master-key")
	n := &Note{
		Payload: "secret payload",
		Draft:   "secret draft",
	}

	if err := n.EncryptNote(encryptor); err != nil {
		t.Fatalf("EncryptNote() error = %v", err)
	}
	if n.Payload == "secret payload" || n.Draft == "secret draft" {
		t.Fatal("fields should be encrypted")
	}

	if err := n.DecryptNote(encryptor); err != nil {
		t.Fatalf("DecryptNote() error = %v", err)
	}
	if n.Payload != "secret payload" {
		t.Errorf("Payload = %q, want %q", n.Payload, "secret payload")
	}
	if n.Draft != "secret draft" {
		t.Errorf("Draft = %q, want %q", n.Draft, "secret draft")
	}
}

func TestNote_EncryptDecrypt_EmptyFields(t *testing.T) {
	encryptor := crypto.NewEncryptor("test-master-key")
	n := &Note{}

	if err := n.EncryptNote(encryptor); err != nil {
		t.Fatalf("EncryptNote() error = %v", err)
	}
	if err := n.DecryptNote(encryptor); err != nil {
		t.Fatalf("DecryptNote() error = %v", err)
	}
}
