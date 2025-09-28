package auth

import (
	"context"
	"mime/multipart"

	"wn/internal/domain/dto/request"
	respDto "wn/internal/domain/dto/response"
	userDto "wn/internal/domain/dto/user"
	"wn/internal/infrastructure/repository/user"
	"wn/pkg/applogger"
	"wn/pkg/trx"

	"github.com/google/uuid"
)

type userService interface {
	CreateUserFromAuthCredentials(ctx context.Context, credintials request.RegisterCredentials) (*userDto.User, error)
	UpdateUser(ctx context.Context, userId uuid.UUID, filter *user.UserUpdateParams) error
	GetUserById(ctx context.Context, userId uuid.UUID, password string) (*userDto.User, error)
}

type fileService interface {
	NewFile(ctx context.Context, file *multipart.FileHeader) (string, error)
}

type Service struct {
	tx     trx.TransactionManager
	logger applogger.Logger

	userService userService
	fileService fileService
}

func NewService(
	tx trx.TransactionManager,
	logger applogger.Logger,
	userService userService,
	fileService fileService,
) *Service {
	return &Service{
		tx:          tx,
		logger:      logger,
		userService: userService,
		fileService: fileService,
	}
}

// todo add reg exp check for password and username and email
func (srv *Service) RegisterUser(ctx context.Context, credentials request.RegisterCredentials) (*respDto.RegisterResponse, error) {
	u, err := srv.userService.CreateUserFromAuthCredentials(ctx, credentials)
	if err != nil {
		return nil, err
	}
	t := true
	err = srv.userService.UpdateUser(ctx, u.Id, &user.UserUpdateParams{
		ConfirmedEmail: &t,
	})
	return &respDto.RegisterResponse{
		UserId: u.Id,
	}, nil

}

func (srv *Service) ChangeProfilePicture(ctx context.Context, req request.ChangeProfilePicture, host string) (*respDto.ChangePictureResponse, error) {
	filename, err := srv.fileService.NewFile(ctx, req.File)
	if err != nil {
		return nil, err
	}
	err = srv.userService.UpdateUser(ctx, req.UserId, &user.UserUpdateParams{
		ImgUrl: &filename,
	})
	return &respDto.ChangePictureResponse{
		NewImgurl: host + "/statics/images/" + filename,
	}, err
}

func (srv *Service) GetUserById(ctx context.Context, userId uuid.UUID, host string) (*userDto.User, error) {
	u, err := srv.userService.GetUserById(ctx, userId, "")
	if err != nil {
		return nil, err
	}
	u.ImgUrl = host + "/statics/images/" + u.ImgUrl
	return u, nil
}
