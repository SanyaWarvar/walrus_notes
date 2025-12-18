package dto

import (
	"wn/internal/entity"

	"github.com/google/uuid"
)

type Note struct {
	Id            uuid.UUID   `json:"id"`
	Title         string      `json:"title"`
	Payload       string      `json:"payload"`
	OwnerId       uuid.UUID   `json:"ownerId"`
	HaveAccess    []uuid.UUID `json:"haveAccess"`
	Position      *Position   `json:"position,omitempty"`
	LinkedWithIn  []uuid.UUID `json:"linkedWithIn,omitempty"`
	LinkedWithOut []uuid.UUID `json:"linkedWithOut,omitempty"`
	Draft         string      `json:"draft"`
	LayoutId      uuid.UUID   `json:"layoutId"`
}

func NotesFromEntities(entities []entity.Note, links []entity.Link) []Note {
	output := make([]Note, 0, len(entities))
	out, in := TransformLinks(links)
	for _, item := range entities {
		output = append(output, Note{
			Id:            item.Id,
			Title:         item.Title,
			Payload:       item.Payload,
			OwnerId:       item.OwnerId,
			HaveAccess:    item.HaveAccess,
			LinkedWithOut: out[item.Id],
			LinkedWithIn:  in[item.Id],
			Draft:         item.Draft,
		})
	}
	return output
}

func NotesFromEntitiesWithPosition(entities []entity.NoteWithPosition, links []entity.Link) []Note {
	output := make([]Note, 0, len(entities))
	out, in := TransformLinks(links)
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
			LinkedWithOut: out[item.Id],
			LinkedWithIn:  in[item.Id],
			Draft:         item.Draft,
			LayoutId:      item.LayoutId,
		})
	}
	return output
}

func TransformLinks(links []entity.Link) (map[uuid.UUID][]uuid.UUID, map[uuid.UUID][]uuid.UUID) {
	out := make(map[uuid.UUID][]uuid.UUID, len(links))
	in := make(map[uuid.UUID][]uuid.UUID, len(links))
	for _, item := range links {
		out[item.FirstNoteId] = append(out[item.FirstNoteId], item.SecondNoteId)
		in[item.SecondNoteId] = append(out[item.SecondNoteId], item.FirstNoteId)
	}
	return out, in
}

type Layout struct {
	Id      uuid.UUID `json:"id"`
	Title   string    `json:"title"`
	OwnerId uuid.UUID `json:"ownerId"`
	IsMain  bool      `json:"isMain"`
	Color   string    `json:"color"`
}

type Position struct {
	XPos float64 `json:"xPos"`
	YPos float64 `json:"yPos"`
}
