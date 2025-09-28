package tokens

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	Id       uuid.UUID `json:"id"`
	UserId   uuid.UUID `json:"userId"`
	AccessId uuid.UUID `json:"accessId"`
	ExpAt    time.Time `json:"expAt"`
}
