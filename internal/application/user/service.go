package auth

import (
	"context"
	"mime/multipart"

	"wn/internal/domain/dto"
	"wn/pkg/applogger"
	"wn/pkg/trx"

	"github.com/google/uuid"
)

type userService interface {
	CreateUserFromAuthCredentials(ctx context.Context, credintials dto.RegisterCredentials) (*dto.User, error)
	UpdateUser(ctx context.Context, userId uuid.UUID, filter *dto.UserUpdateParams) error
	GetUserById(ctx context.Context, userId uuid.UUID, password string) (*dto.User, error)
}

type fileService interface {
	NewFile(ctx context.Context, file *multipart.FileHeader) (string, error)
}

type layoutService interface {
	CreateLayout(ctx context.Context, title, color string, ownerId uuid.UUID, isMain bool) (uuid.UUID, error)
}

type Service struct {
	tx     trx.TransactionManager
	logger applogger.Logger

	userService   userService
	fileService   fileService
	layoutService layoutService
}

func NewService(
	tx trx.TransactionManager,
	logger applogger.Logger,
	userService userService,
	fileService fileService,
	layoutService layoutService,
) *Service {
	return &Service{
		tx:            tx,
		logger:        logger,
		userService:   userService,
		fileService:   fileService,
		layoutService: layoutService,
	}
}

// todo add reg exp check for password and username and email
func (srv *Service) RegisterUser(ctx context.Context, credentials dto.RegisterCredentials) (*dto.RegisterResponse, error) {
	var u *dto.User
	var err error
	if err = srv.tx.Transaction(ctx, func(ctx context.Context) error {
		u, err = srv.userService.CreateUserFromAuthCredentials(ctx, credentials)
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
		_, err = srv.layoutService.CreateLayout(ctx, "All Notes", "#FFFFFF", u.Id, true)
		return err
	}); err != nil {
		return nil, err
	}

	return &dto.RegisterResponse{
		UserId: u.Id,
	}, nil
}

func (srv *Service) ChangeProfilePicture(ctx context.Context, req dto.ChangeProfilePicture, host string) (*dto.ChangePictureResponse, error) {
	filename, err := srv.fileService.NewFile(ctx, req.File)
	if err != nil {
		return nil, err
	}
	err = srv.userService.UpdateUser(ctx, req.UserId, &dto.UserUpdateParams{
		ImgUrl: &filename,
	})
	return &dto.ChangePictureResponse{
		NewImgurl: host + "/statics/images/" + filename,
	}, err
}

func (srv *Service) GetUserById(ctx context.Context, userId uuid.UUID, host string) (*dto.User, error) {
	u, err := srv.userService.GetUserById(ctx, userId, "")
	if err != nil {
		return nil, err
	}
	u.ImgUrl = host + "/statics/images/" + u.ImgUrl
	return u, nil
}
