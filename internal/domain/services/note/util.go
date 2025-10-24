package note

import (
	"wn/internal/entity"

	"github.com/google/uuid"
)

func getIds(notes []entity.NoteWithPosition) []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(notes))
	for _, n := range notes {
		ids = append(ids, n.Id)
	}
	return ids
}
