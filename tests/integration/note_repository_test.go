//go:build integration

package integration

import (
	"testing"
	"wn/internal/domain/entity"
	"wn/internal/domain/services/crypto"
	layoutrepo "wn/internal/infrastructure/repository/layout"
	noterepo "wn/internal/infrastructure/repository/note"
	"wn/pkg/util"

	"github.com/google/uuid"
)

func TestNoteRepository_CreateAndGet(t *testing.T) {
	resetDB(t)
	user := seedUser(t)
	layoutRepo := layoutrepo.NewRepository(env.DB)
	noteRepo := noterepo.NewRepository(env.DB)

	layout := &entity.Layout{
		Id:         uuid.New(),
		Title:      "Notes",
		OwnerId:    user.Id,
		HaveAccess: []uuid.UUID{user.Id},
		IsMain:     true,
		Color:      "#FFFFFF",
	}
	if _, err := layoutRepo.CreateLayout(env.Ctx, layout); err != nil {
		t.Fatalf("CreateLayout() error = %v", err)
	}

	encryptor := crypto.NewEncryptor("integration-encrypt-key")
	note := &entity.Note{
		Id:         uuid.New(),
		Title:      "Test Note",
		Payload:    "secret content",
		CreatedAt:  util.GetCurrentUTCTime(),
		OwnerId:    user.Id,
		HaveAccess: []uuid.UUID{user.Id},
		LayoutId:   layout.Id,
	}
	if err := note.EncryptNote(encryptor); err != nil {
		t.Fatalf("EncryptNote() error = %v", err)
	}

	id, err := noteRepo.CreateNote(env.Ctx, note)
	if err != nil {
		t.Fatalf("CreateNote() error = %v", err)
	}
	if id != note.Id {
		t.Errorf("id = %v, want %v", id, note.Id)
	}

	got, err := noteRepo.GetById(env.Ctx, note.Id)
	if err != nil {
		t.Fatalf("GetById() error = %v", err)
	}
	if err := got.DecryptNote(encryptor); err != nil {
		t.Fatalf("DecryptNote() error = %v", err)
	}
	if got.Payload != "secret content" {
		t.Errorf("Payload = %q, want secret content", got.Payload)
	}
}
