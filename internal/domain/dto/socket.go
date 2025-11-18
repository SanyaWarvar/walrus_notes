package dto

import (
	"encoding/json"

	"github.com/google/uuid"
)

type SocketMessage struct {
	Event   string          `json:"event"`
	Payload json.RawMessage `json:"payload"`
}

type DraftNote struct {
	NoteId   uuid.UUID `json:"noteId"`
	NewDraft string    `json:"newDraft"`
}
