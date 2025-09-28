package auth

import (
	"wn/internal/domain/enum"
	"time"
)

type ConfirmationCode struct {
	Code      string               `json:"code"`
	CreatedAt time.Time            `json:"createdAt"`
	Action    enum.EmailCodeAction `json:"action"`
}
