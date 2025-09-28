package dto

import (
	"time"

	"github.com/google/uuid"
)

type RegisterResponse struct {
	UserId uuid.UUID `json:"userId"`
}

type SendCodeResponse struct {
	NextCodeDelay time.Duration `json:"nextCodeDelay"`
}

type ChangePictureResponse struct {
	NewImgurl string `json:"newImgUrl"`
}
