package entity

import (
	"time"

	"github.com/google/uuid"
)

type Note struct {
	Id         uuid.UUID   `json:"id" db:"id"`
	Title      string      `json:"title" db:"title"`
	Payload    string      `json:"payload" db:"payload"`
	CreatedAt  time.Time   `json:"createdAt" db:"created_at"`
	OwnerId    uuid.UUID   `json:"ownerId" db:"owner_id"`
	HaveAccess []uuid.UUID `json:"haveAccess" db:"have_access"`
	Draft      string      `json:"draft" db:"draft"`
}

func (n Note) GetId() uuid.UUID {
	return n.Id
}

type NoteWithPosition struct {
	Note         `json:"note"`
	NotePosition `json:"notePosition"`
}

func (n NoteWithPosition) GetId() uuid.UUID {
	return n.Id
}

type NotePosition struct {
	NoteId    uuid.UUID `json:"noteId" db:"note_id"`
	LayoutId  uuid.UUID `json:"layoutId" db:"layout_id"`
	XPosition float64   `json:"xPosition" db:"x_position"`
	YPosition float64   `json:"yPosition" db:"y_position"`
}

type Layout struct {
	Id         uuid.UUID   `json:"id" db:"id"`
	Title      string      `json:"title"`
	OwnerId    uuid.UUID   `json:"ownerId" db:"owner_id"`
	HaveAccess []uuid.UUID `json:"haveAccess" db:"have_access"`
	IsMain     bool        `json:"isMain"`
}

type Link struct {
	LayoutId     uuid.UUID `json:"layoutId"`
	FirstNoteId  uuid.UUID `json:"firstNoteId"`
	SecondNoteId uuid.UUID `json:"secondNoteId"`
}
