package events

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func assertSocketEvent(t *testing.T, eventName string, payload any, gotEvent string, gotPayload json.RawMessage) {
	t.Helper()

	if gotEvent != eventName {
		t.Errorf("Event = %q, want %q", gotEvent, eventName)
	}

	expected, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if string(gotPayload) != string(expected) {
		t.Errorf("Payload = %s, want %s", gotPayload, expected)
	}
}

func TestDeleteLayoutEvent_ToSocketEvent(t *testing.T) {
	layoutID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	e := &DeleteLayoutEvent{LayoutId: layoutID}
	msg := e.ToSocketEvent()
	assertSocketEvent(t, "DELETE_LAYOUT", e, msg.Event, msg.Payload)
}

func TestUpdateLayoutEvent_ToSocketEvent(t *testing.T) {
	layoutID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	e := &UpdateLayoutEvent{LayoutId: layoutID}
	msg := e.ToSocketEvent()
	assertSocketEvent(t, "UPDATE_LAYOUT", e, msg.Event, msg.Payload)
}

func TestCreateNoteEvent_ToSocketEvent(t *testing.T) {
	layoutID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	noteID := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	e := &CreateNoteEvent{LayoutId: layoutID, NoteId: noteID}
	msg := e.ToSocketEvent()
	assertSocketEvent(t, "CREATE_NOTE", e, msg.Event, msg.Payload)
}

func TestUpdateNoteEvent_ToSocketEvent(t *testing.T) {
	noteID := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	e := &UpdateNoteEvent{NoteId: noteID}
	msg := e.ToSocketEvent()
	assertSocketEvent(t, "UPDATE_NOTE", e, msg.Event, msg.Payload)
}

func TestDeleteNoteEvent_ToSocketEvent(t *testing.T) {
	noteID := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	e := &DeleteNoteEvent{NoteId: noteID}
	msg := e.ToSocketEvent()
	assertSocketEvent(t, "DELETE_NOTE", e, msg.Event, msg.Payload)
}

func TestDragNoteEvent_ToSocketEvent(t *testing.T) {
	noteID := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	layoutID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	e := &DragNoteEvent{NoteId: noteID, ToLayoutId: layoutID}
	msg := e.ToSocketEvent()
	assertSocketEvent(t, "DRAG_NOTE", e, msg.Event, msg.Payload)
}

func TestChangeNotePositionEvent_ToSocketEvent(t *testing.T) {
	layoutID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	e := &ChangeNotePositionEvent{LayoutId: layoutID}
	msg := e.ToSocketEvent()
	assertSocketEvent(t, "CHANGE_NOTE_POSITION", e, msg.Event, msg.Payload)
}

func TestChangeLinkEvent_ToSocketEvent(t *testing.T) {
	layoutID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	e := &ChangeLinkEvent{LayoutId: layoutID}
	msg := e.ToSocketEvent()
	assertSocketEvent(t, "CHANGE_NOTE_LINKS", e, msg.Event, msg.Payload)
}
