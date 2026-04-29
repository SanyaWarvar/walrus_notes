package dto

import "github.com/google/uuid"

type ConnectionID string

type Connection interface {
	ID() ConnectionID
	UserID() uuid.UUID
	Send(msg *SocketMessage) error
	SendPing() error
	Close() error
	ReadMessage() (*SocketMessage, error)
}
