package events

import (
	"encoding/json"
	"wn/internal/domain/dto"

	"github.com/google/uuid"
)

type DeleteLayoutEvent struct {
	LayoutId uuid.UUID
}

func (e *DeleteLayoutEvent) ToSocketEvent() *dto.SocketMessage {
	d, _ := json.Marshal(e)
	return &dto.SocketMessage{
		Event:   "DELETE_LAYOUT",
		Payload: d,
	}
}
