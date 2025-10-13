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
}

type NoteWithPosition struct {
	Note         `json:"note"`
	NotePosition `json:"notePosition"`
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
}

type Link struct {
	Id         uuid.UUID `json:"id" db:"id"`
	X1Position float64   `json:"x1Position" db:"x1_position"`
	Y1Position float64   `json:"y1Position" db:"y1_position"`
	X2Position float64   `json:"x2Position" db:"x2_position"`
	Y2Position float64   `json:"y2Position" db:"y2_position"`
	Color      string    `json:"color" db:"color"`
}
