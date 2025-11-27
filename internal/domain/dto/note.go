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
	LinkedWith []uuid.UUID `json:"linkedWith,omitempty"`
	Draft      string      `json:"draft"`
	LayoutId   uuid.UUID   `json:"layoutId"`
}

func NotesFromEntities(entities []entity.Note, links []entity.Link) []Note {
	output := make([]Note, 0, len(entities))
	transformedLinks := TransformLinks(links)
	for _, item := range entities {
		output = append(output, Note{
			Id:         item.Id,
			Title:      item.Title,
			Payload:    item.Payload,
			OwnerId:    item.OwnerId,
			HaveAccess: item.HaveAccess,
			LinkedWith: transformedLinks[item.Id],
			Draft:      item.Draft,
		})
	}
	return output
}

func NotesFromEntitiesWithPosition(entities []entity.NoteWithPosition, links []entity.Link) []Note {
	output := make([]Note, 0, len(entities))
	transformedLinks := TransformLinks(links)
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
			LinkedWith: transformedLinks[item.Id],
			Draft:      item.Draft,
		})
	}
	return output
}

func TransformLinks(links []entity.Link) map[uuid.UUID][]uuid.UUID {
	output := make(map[uuid.UUID][]uuid.UUID, len(links))
	for _, item := range links {
		output[item.FirstNoteId] = append(output[item.FirstNoteId], item.SecondNoteId)
	}
	return output
}

type Layout struct {
	Id      uuid.UUID `json:"id"`
	Title   string    `json:"title"`
	OwnerId uuid.UUID `json:"ownerId"`
	IsMain  bool      `json:"isMain"`
}

type Position struct {
	XPos float64 `json:"xPos,omitempty"`
	YPos float64 `json:"yPos,omitempty"`
}
