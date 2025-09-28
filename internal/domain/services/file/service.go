package file

import (
	"context"
	"encoding/base64"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"wn/internal/infrastructure/repository/file"
	"wn/pkg/applogger"
	"wn/pkg/util"

	"github.com/pkg/errors"
)

const filePath = "./statics/images/"

type fileRepo interface {
	GetAllFiles(ctx context.Context) ([]file.StaticFile, error)
	CreateFile(ctx context.Context, filename string, encodedFile string) error
}

type Service struct {
	logger applogger.Logger

	fileRepo fileRepo
}

func NewService(lgr applogger.Logger, fileRepo fileRepo) *Service {
	return &Service{
		logger:   lgr,
		fileRepo: fileRepo,
	}
}

func (srv *Service) GenerateStatics(ctx context.Context) error {
	files, err := srv.fileRepo.GetAllFiles(ctx)
	if err != nil {
		return errors.Wrap(err, "srv.fileRepo.GetAllFiles")
	}

	for ind, item := range files {
		if item.FileAsString == "" {
			srv.logger.WithCtx(ctx).Debugf("Skipping empty file: %s", item.Filename)
			continue
		}

		// Декодируем base64
		files[ind].File, err = base64.RawStdEncoding.DecodeString(item.FileAsString)
		if err != nil {
			srv.logger.WithCtx(ctx).Warnf("Base64 decode failed for %s: %s", item.Filename, err.Error())
			continue
		}

		fullPath := filepath.Join(filePath, item.Filename)

		// Проверяем существование файла
		if _, err := os.Stat(fullPath); err == nil {
			// Файл уже существует, пропускаем
			srv.logger.WithCtx(ctx).Debugf("File already exists, skipping: %s", item.Filename)
			continue
		} else if !os.IsNotExist(err) {
			// Другая ошибка при проверке файла
			srv.logger.WithCtx(ctx).Warnf("File stat error for %s: %s", item.Filename, err.Error())
			continue
		}

		// Атомарная запись файла
		tempPath := fullPath + ".tmp"
		if err := os.WriteFile(tempPath, files[ind].File, 0644); err != nil {
			srv.logger.WithCtx(ctx).Errorf("Failed to write file %s: %s", item.Filename, err.Error())
			continue
		}

		// Переименовываем временный файл в целевой
		if err := os.Rename(tempPath, fullPath); err != nil {
			srv.logger.WithCtx(ctx).Errorf("Failed to rename temp file %s: %s", item.Filename, err.Error())
			continue
		}

		srv.logger.WithCtx(ctx).Infof("Successfully generated file: %s", item.Filename)
	}
	return nil
}

func (srv *Service) NewFile(ctx context.Context, file *multipart.FileHeader) (string, error) {
	ext := strings.TrimPrefix(file.Filename, ".")
	filename := util.NewUUID().String() + "." + ext
	encodedFile, err := encodeFile(file)
	if err != nil {
		return "", err
	}
	err = srv.fileRepo.CreateFile(ctx, filename, encodedFile)
	return filename, err
}

func encodeFile(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()
	fileData, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	encodedFile := base64.RawStdEncoding.EncodeToString(fileData)
	return encodedFile, nil
}
