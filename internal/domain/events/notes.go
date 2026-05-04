package events

import (
	"encoding/json"
	"wn/internal/domain/dto"

	"github.com/google/uuid"
)

type CreateNoteEvent struct {
	LayoutId uuid.UUID `json:"layoutId"`
	NoteId   uuid.UUID `json:"noteId"`
}

func (e *CreateNoteEvent) ToSocketEvent() *dto.SocketMessage {
	d, _ := json.Marshal(e)
	return &dto.SocketMessage{
		Event:   "CREATE_NOTE",
		Payload: d,
	}
}

type UpdateNoteEvent struct {
	NoteId uuid.UUID `json:"noteId"`
}

func (e *UpdateNoteEvent) ToSocketEvent() *dto.SocketMessage {
	d, _ := json.Marshal(e)
	return &dto.SocketMessage{
		Event:   "UPDATE_NOTE",
		Payload: d,
	}
}

type DeleteNoteEvent struct {
	NoteId uuid.UUID `json:"noteId"`
}

func (e *DeleteNoteEvent) ToSocketEvent() *dto.SocketMessage {
	d, _ := json.Marshal(e)
	return &dto.SocketMessage{
		Event:   "DELETE_NOTE",
		Payload: d,
	}
}

type DragNoteEvent struct {
	NoteId     uuid.UUID `json:"noteId"`
	ToLayoutId uuid.UUID `json:"toLayoutId"`
}

func (e *DragNoteEvent) ToSocketEvent() *dto.SocketMessage {
	d, _ := json.Marshal(e)
	return &dto.SocketMessage{
		Event:   "DRAG_NOTE",
		Payload: d,
	}
}

type ChangeNotePositionEvent struct {
	LayoutId uuid.UUID `json:"layoutId"`
}

func (e *ChangeNotePositionEvent) ToSocketEvent() *dto.SocketMessage {
	d, _ := json.Marshal(e)
	return &dto.SocketMessage{
		Event:   "CHANGE_NOTE_POSITION",
		Payload: d,
	}
}

type ChangeLinkEvent struct {
	LayoutId uuid.UUID `json:"layoutId"`
}

func (e *ChangeLinkEvent) ToSocketEvent() *dto.SocketMessage {
	d, _ := json.Marshal(e)
	return &dto.SocketMessage{
		Event:   "CHANGE_NOTE_LINKS",
		Payload: d,
	}
}
