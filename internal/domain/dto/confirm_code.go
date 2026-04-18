package dto

import (
	"time"
	"wn/internal/domain/enum"
)

type ConfirmationCode struct {
	Code      string               `json:"code"`
	CreatedAt time.Time            `json:"createdAt"`
	Action    enum.EmailCodeAction `json:"action"`
}
