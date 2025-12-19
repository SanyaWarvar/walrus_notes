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

	for _, item := range files {
		createFile(&item)
	}
	return nil
}

func (srv *Service) NewFile(ctx context.Context, fileObj *multipart.FileHeader) (string, error) {
	ext := strings.TrimPrefix(fileObj.Filename, ".")
	filename := util.NewUUID().String() + "." + ext
	encodedFile, err := encodeFile(fileObj)
	if err != nil {
		return "", err
	}
	err = srv.fileRepo.CreateFile(ctx, filename, encodedFile)
	if err != nil {
		return filename, err
	}
	f := &file.StaticFile{
		Filename:     filename,
		FileAsString: encodedFile,
		File:         []byte{},
	}
	return filename, createFile(f)
}

func createFile(item *file.StaticFile) error {
	var err error
	if item.FileAsString == "" {
		return errors.Wrap(err, "item.FileAsString")
	}

	item.File, err = base64.RawStdEncoding.DecodeString(item.FileAsString)
	if err != nil {
		return errors.Wrap(err, "base64.RawStdEncoding.DecodeString")
	}

	fullPath := filepath.Join(filePath, item.Filename)

	// Проверяем существование файла
	if _, err := os.Stat(fullPath); err == nil {

	} else if !os.IsNotExist(err) {
		return errors.Wrap(err, "os.IsNotExist")
	}

	// Атомарная запись файла
	tempPath := fullPath + ".tmp"
	if err := os.WriteFile(tempPath, item.File, 0644); err != nil {
		return errors.Wrap(err, "os.WriteFile")
	}

	// Переименовываем временный файл в целевой
	if err := os.Rename(tempPath, fullPath); err != nil {
		return errors.Wrap(err, "os.Rename(tempPath, fullPath)")
	}

	return nil
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
