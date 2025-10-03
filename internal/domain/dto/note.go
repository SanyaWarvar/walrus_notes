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

type Layout struct {
	Id      uuid.UUID `json:"id"`
	Title   string    `json:"title"`
	OwnerId uuid.UUID `json:"ownerId"`
}
