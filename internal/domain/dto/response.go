package dto

import (
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
