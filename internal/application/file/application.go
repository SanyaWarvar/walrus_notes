package file

import (
	"context"
	"mime/multipart"
	"wn/internal/domain/dto"
	"wn/internal/domain/dto/request"
	"wn/pkg/applogger"
	"wn/pkg/trx"

	"github.com/google/uuid"
)

type fileService interface {
	NewFile(ctx context.Context, file *multipart.FileHeader) (string, error)
}

type Service struct {
	tx     trx.TransactionManager
	logger applogger.Logger

	fileService fileService
}

func NewService(
	tx trx.TransactionManager,
	logger applogger.Logger,
	fileService fileService,
) *Service {
	return &Service{
		tx:          tx,
		logger:      logger,
		fileService: fileService,
	}
}

func (srv *Service) UploadFile(ctx context.Context, userId uuid.UUID, req request.UploadFileRequest, host string) (*dto.UploadFileResponse, error) {
	filename, err := srv.fileService.NewFile(ctx, req.File)
	if err != nil {
		return nil, err
	}
	return &dto.UploadFileResponse{
		ImgUrl: host + "/statics/images/" + filename,
	}, err
}
