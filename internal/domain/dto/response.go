package dto

import (
	"time"

	"github.com/google/uuid"
)

type UploadFileResponse struct {
	ImgUrl string `json:"imgUrl"`
}

type ExportInfoRequest struct {
	UserId uuid.UUID `json:"userId"`
}

type ImportInfoRequest struct {
	Info ExportInfo `json:"info"`
}

type RegisterResponse struct {
	UserId uuid.UUID `json:"userId"`
}

type SendCodeResponse struct {
	NextCodeDelay time.Duration `json:"nextCodeDelay"`
}

type ChangePictureResponse struct {
	NewImgurl string `json:"newImgUrl"`
}

type LoginResponse struct {
	Access  string    `json:"accessToken"`
	Refresh string    `json:"refreshToken"`
	UserId  uuid.UUID `json:"userId"`
}
