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

type NotePosition struct {
	NoteId    uuid.UUID `json:"noteId" db:"note_id"`
	LayoutIdd uuid.UUID `json:"layoutId" db:"layout_id"`
	XPosition float64   `json:"xPosition" db:"x_position"`
	YPosition float64   `json:"yPosition" db:"y_position"`
}

type Layout struct {
	LayoutID   uuid.UUID   `json:"layoutId" db:"layout_id"`
	OwnerID    uuid.UUID   `json:"ownerId" db:"owner_id"`
	HaveAccess []uuid.UUID `json:"haveAccess" db:"have_access"`
}

type Link struct {
	Id        uuid.UUID `json:"id" db:"id"`
	XPosition float64   `json:"xPosition" db:"x_position"`
	YPosition float64   `json:"yPosition" db:"y_position"`
	Color     string    `json:"color" db:"color"`
}
