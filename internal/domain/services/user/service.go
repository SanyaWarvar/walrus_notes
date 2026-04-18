package auth

import (
	"context"
	"crypto/sha512"
	"fmt"
	"wn/internal/domain/dto"
	"wn/internal/domain/entity"
	apperrors "wn/internal/errors"

	"wn/pkg/applogger"
	"wn/pkg/constants"
	"wn/pkg/trx"
	"wn/pkg/util"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type userRepo interface {
	CreateUser(ctx context.Context, item *entity.User) error
	GetUser(ctx context.Context, filter dto.UserFilter) (*entity.User, bool, error)
	UpdateUser(ctx context.Context, userId uuid.UUID, updateParams *dto.UserUpdateParams) error
}

type Service struct {
	passwordSalt string

	tx     trx.TransactionManager
	logger applogger.Logger

	userRepo userRepo
}

func NewService(
	passwordSalt string,
	tx trx.TransactionManager,
	logger applogger.Logger,
	userRepo userRepo,
) *Service {
	return &Service{
		passwordSalt: passwordSalt,
		tx:           tx,
		logger:       logger,
		userRepo:     userRepo,
	}
}

func (srv *Service) CreateUserFromAuthCredentials(ctx context.Context, credintials dto.RegisterCredentials) (*dto.User, error) {
	user := dto.User{
		Id:        util.NewUUID(),
		Username:  credintials.Username,
		Email:     credintials.Email,
		ImgUrl:    "base.png",
		CreatedAt: util.GetCurrentUTCTime(),
	}
	userEntity := entity.User{
		Id:             user.Id,
		Username:       user.Username,
		Email:          user.Email,
		ConfirmedEmail: false,
		Password:       srv.generatePasswordHash(credintials.Password),
		ImgUrl:         "base.png",
		CreatedAt:      user.CreatedAt,
		Role:           constants.ClientRole,
	}
	err := srv.userRepo.CreateUser(ctx, &userEntity)
	return &user, err
}

func (srv *Service) UpdateUser(ctx context.Context, userId uuid.UUID, filter *dto.UserUpdateParams) error {
	_, ex, err := srv.userRepo.GetUser(ctx, dto.UserFilter{
		Id: &userId,
	})
	if err != nil {
		return errors.Wrap(err, "srv.userRepo.GetUser")
	}

	if !ex {
		return apperrors.UserNotFound
	}

	if filter.Password != nil {
		newPassword := srv.generatePasswordHash(*filter.Password)
		filter.Password = &newPassword
	}

	return srv.userRepo.UpdateUser(ctx, userId, filter)
}

func (srv *Service) GetUserById(ctx context.Context, userId uuid.UUID, password string) (*dto.User, error) {
	targetEntityUser, ex, err := srv.userRepo.GetUser(ctx, dto.UserFilter{
		Id: &userId,
	})
	if err != nil {
		return nil, err
	}

	if !ex {
		return nil, apperrors.UserNotFound
	}

	if password != "" {
		if !srv.comparePassword(password, targetEntityUser.Password) {
			return nil, apperrors.IncorrectPassword
		}
	}
	return dto.UserDtoFromEntity(targetEntityUser), nil
}

func (srv *Service) GetUserByEmail(ctx context.Context, email string, password string) (*dto.User, error) {
	targetEntityUser, ex, err := srv.userRepo.GetUser(ctx, dto.UserFilter{
		Email: &email,
	})
	if err != nil {
		return nil, err
	}

	if !ex {
		return nil, apperrors.UserNotFound
	}

	if password != "" {
		if !srv.comparePassword(password, targetEntityUser.Password) {
			return nil, apperrors.IncorrectPassword
		}
	}
	return dto.UserDtoFromEntity(targetEntityUser), nil
}

func (srv *Service) comparePassword(origin, existed string) bool {
	return srv.generatePasswordHash(origin) == existed
}

func (srv *Service) generatePasswordHash(password string) string {
	hash := sha512.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(srv.passwordSalt)))
}
