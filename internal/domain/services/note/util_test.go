package note

import (
	"testing"
	"wn/internal/domain/entity"

	"github.com/google/uuid"
)

func TestGetIds(t *testing.T) {
	id1 := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	id2 := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

	notes := []entity.Note{
		{Id: id1},
		{Id: id2},
	}

	got := getIds(notes)
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0] != id1 || got[1] != id2 {
		t.Errorf("getIds() = %v, want [%v, %v]", got, id1, id2)
	}
}

func TestGetIds_Empty(t *testing.T) {
	got := getIds([]entity.Note{})
	if len(got) != 0 {
		t.Errorf("len = %d, want 0", len(got))
	}
}
