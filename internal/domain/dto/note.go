package dto

import "github.com/google/uuid"

type Note struct {
	Id         uuid.UUID   `json:"id"`
	Title      string      `json:"title"`
	Payload    string      `json:"payload"`
	OwnerId    uuid.UUID   `json:"ownerId"`
	HaveAccess []uuid.UUID `json:"haveAccess"`
}
