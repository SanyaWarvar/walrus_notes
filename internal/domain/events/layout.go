package events

import (
	"encoding/json"
	"wn/internal/domain/dto"

	"github.com/google/uuid"
)

type DeleteLayoutEvent struct {
	LayoutId uuid.UUID `json:"layoutId"`
}

func (e *DeleteLayoutEvent) ToSocketEvent() *dto.SocketMessage {
	d, _ := json.Marshal(e)
	return &dto.SocketMessage{
		Event:   "DELETE_LAYOUT",
		Payload: d,
	}
}

type UpdateLayoutEvent struct {
	LayoutId uuid.UUID `json:"layoutId"`
}

func (e *UpdateLayoutEvent) ToSocketEvent() *dto.SocketMessage {
	d, _ := json.Marshal(e)
	return &dto.SocketMessage{
		Event:   "UPDATE_LAYOUT",
		Payload: d,
	}
}
