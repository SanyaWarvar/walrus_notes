package events

import "wn/internal/domain/dto"

type Event interface {
	ToSocketEvent() *dto.SocketMessage
}
