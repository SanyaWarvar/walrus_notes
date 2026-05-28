package dto

import (
	"testing"
	"wn/internal/domain/entity"

	"github.com/google/uuid"
)

func TestTransformLinks(t *testing.T) {
	n1 := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	n2 := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	n3 := uuid.MustParse("6ba7b811-9dad-11d1-80b4-00c04fd430c8")

	links := []entity.Link{
		{FirstNoteId: n1, SecondNoteId: n2},
		{FirstNoteId: n1, SecondNoteId: n3},
	}

	out, in := TransformLinks(links)

	if len(out[n1]) != 2 {
		t.Fatalf("out[n1] len = %d, want 2", len(out[n1]))
	}
	if len(in[n2]) != 1 || in[n2][0] != n1 {
		t.Errorf("in[n2] = %v, want [%v]", in[n2], n1)
	}
}

func TestNotesFromEntities(t *testing.T) {
	n1 := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	n2 := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	owner := uuid.MustParse("6ba7b812-9dad-11d1-80b4-00c04fd430c8")
	layout := uuid.MustParse("6ba7b813-9dad-11d1-80b4-00c04fd430c8")

	entities := []entity.Note{
		{Id: n1, Title: "A", OwnerId: owner, LayoutId: layout},
		{Id: n2, Title: "B", OwnerId: owner, LayoutId: layout},
	}
	links := []entity.Link{{FirstNoteId: n1, SecondNoteId: n2}}

	got := NotesFromEntities(entities, links)
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if len(got[0].LinkedWithOut) != 1 || got[0].LinkedWithOut[0] != n2 {
		t.Errorf("LinkedWithOut = %v, want [%v]", got[0].LinkedWithOut, n2)
	}
}

func TestNotesFromEntitiesWithPosition(t *testing.T) {
	n1 := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	owner := uuid.MustParse("6ba7b812-9dad-11d1-80b4-00c04fd430c8")
	layout := uuid.MustParse("6ba7b813-9dad-11d1-80b4-00c04fd430c8")

	entities := []entity.NoteWithPosition{
		{
			Note:         entity.Note{Id: n1, Title: "A", OwnerId: owner, LayoutId: layout},
			NotePosition: entity.NotePosition{XPosition: 10, YPosition: 20},
		},
	}

	got := NotesFromEntitiesWithPosition(entities, nil)
	if got[0].Position == nil {
		t.Fatal("Position should be set")
	}
	if got[0].Position.XPos != 10 || got[0].Position.YPos != 20 {
		t.Errorf("Position = %+v, want (10, 20)", got[0].Position)
	}
}
