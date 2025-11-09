package note

import (
	"github.com/google/uuid"
)

type ItemWithId interface {
	GetId() uuid.UUID
}

func getIds[T ItemWithId](notes []T) []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(notes))
	for _, n := range notes {
		ids = append(ids, n.GetId())
	}
	return ids
}
