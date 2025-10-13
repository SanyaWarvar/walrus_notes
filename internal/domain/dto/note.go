package dto

import (
	"wn/internal/entity"

	"github.com/google/uuid"
)

type Note struct {
	Id         uuid.UUID   `json:"id"`
	Title      string      `json:"title"`
	Payload    string      `json:"payload"`
	OwnerId    uuid.UUID   `json:"ownerId"`
	HaveAccess []uuid.UUID `json:"haveAccess"`
	Position   *Position   `json:"position,omitempty"`
}

func NotesFromEntities(entities []entity.Note) []Note {
	output := make([]Note, 0, len(entities))
	for _, item := range entities {
		output = append(output, Note{
			Id:         item.Id,
			Title:      item.Title,
			Payload:    item.Payload,
			OwnerId:    item.OwnerId,
			HaveAccess: item.HaveAccess,
		})
	}
	return output
}

func NotesFromEntitiesWithPosition(entities []entity.NoteWithPosition) []Note {
	output := make([]Note, 0, len(entities))
	for _, item := range entities {
		output = append(output, Note{
			Id:         item.Id,
			Title:      item.Title,
			Payload:    item.Payload,
			OwnerId:    item.OwnerId,
			HaveAccess: item.HaveAccess,
			Position: &Position{
				XPos: item.XPosition,
				YPos: item.YPosition,
			},
		})
	}
	return output
}

type Layout struct {
	Id      uuid.UUID `json:"id"`
	Title   string    `json:"title"`
	OwnerId uuid.UUID `json:"ownerId"`
}

type Position struct {
	XPos float64 `json:"xPos,omitempty"`
	YPos float64 `json:"yPos,omitempty"`
}
