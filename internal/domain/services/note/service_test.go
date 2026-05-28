package note

import (
	"testing"
	"wn/internal/domain/dto"

	"github.com/google/uuid"
)

func TestGenerateCluster_Empty(t *testing.T) {
	srv := &Service{}
	got := srv.GenerateCluster(nil)
	if got != nil {
		t.Errorf("GenerateCluster(nil) = %v, want nil", got)
	}
}

func TestGenerateCluster_OffsetsNotes(t *testing.T) {
	srv := &Service{}
	layout1 := uuid.New()
	layout2 := uuid.New()

	notes := []dto.Note{
		{
			LayoutId: layout1,
			Position: &dto.Position{XPos: 10, YPos: 10},
		},
		{
			LayoutId: layout1,
			Position: &dto.Position{XPos: 50, YPos: 50},
		},
		{
			LayoutId: layout2,
			Position: &dto.Position{XPos: 100, YPos: 100},
		},
	}

	got := srv.GenerateCluster(notes)
	if len(got) != 3 {
		t.Fatalf("len = %d, want 3", len(got))
	}

	// Кластеры должны быть разнесены — координаты не совпадают с исходными для всех заметок сразу
	sameAsInput := 0
	for i, n := range got {
		if n.Position == nil {
			t.Fatalf("note %d has nil position", i)
		}
		if n.Position.XPos == notes[i].Position.XPos && n.Position.YPos == notes[i].Position.YPos {
			sameAsInput++
		}
	}
	if sameAsInput == len(got) {
		t.Error("expected cluster offsets to change note positions")
	}
}

func TestGenerateCluster_NotesWithoutPosition(t *testing.T) {
	srv := &Service{}
	layout := uuid.New()

	got := srv.GenerateCluster([]dto.Note{
		{LayoutId: layout, Title: "no pos"},
	})
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
}
