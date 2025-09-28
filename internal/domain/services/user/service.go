package auth

import (
	"context"
	"crypto/sha512"
	"fmt"
	"wn/internal/domain/dto/request"
	"wn/internal/domain/dto/user"
	apperrors "wn/internal/errors"
	userRepository "wn/internal/infrastructure/repository/user"

	"wn/pkg/applogger"
	"wn/pkg/constants"
	"wn/pkg/trx"
	"wn/pkg/util"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// todo вынести в конфиг
var salt string = "pashatechnik"

type userRepo interface {
	CreateUser(ctx context.Context, item *userRepository.User) error
	GetUser(ctx context.Context, filter userRepository.UserFilter) (*userRepository.User, bool, error)
	UpdateUser(ctx context.Context, userId uuid.UUID, updateParams *userRepository.UserUpdateParams) error
}

type Service struct {
	tx     trx.TransactionManager
	logger applogger.Logger

	userRepo userRepo
}

func NewService(
	tx trx.TransactionManager,
	logger applogger.Logger,
	userRepo userRepo,
) *Service {
	return &Service{
		tx:       tx,
		logger:   logger,
		userRepo: userRepo,
	}
}

func (srv *Service) CreateUserFromAuthCredentials(ctx context.Context, credintials request.RegisterCredentials) (*user.User, error) {
	user := user.User{
		Id:        util.NewUUID(),
		Username:  credintials.Username,
		Email:     credintials.Email,
		ImgUrl:    "base.png",
		CreatedAt: util.GetCurrentUTCTime(),
	}
	userEntity := userRepository.User{
		Id:        user.Id,
		Username:  user.Username,
		Email:     user.Email,
		Password:  generatePasswordHash(credintials.Password),
		ImgUrl:    "base.png",
		CreatedAt: user.CreatedAt,
		Role:      constants.ClientRole,
	}
	err := srv.userRepo.CreateUser(ctx, &userEntity)
	return &user, err
}

func (srv *Service) UpdateUser(ctx context.Context, userId uuid.UUID, filter *userRepository.UserUpdateParams) error {
	_, ex, err := srv.userRepo.GetUser(ctx, userRepository.UserFilter{
		Id: &userId,
	})
	if err != nil {
		return errors.Wrap(err, "srv.userRepo.GetUser")
	}

	if !ex {
		return apperrors.UserNotFound
	}

	if filter.Password != nil {
		newPassword := generatePasswordHash(*filter.Password)
		filter.Password = &newPassword
	}

	return srv.userRepo.UpdateUser(ctx, userId, filter)
}

func (srv *Service) GetUserById(ctx context.Context, userId uuid.UUID, password string) (*user.User, error) {
	targetEntityUser, ex, err := srv.userRepo.GetUser(ctx, userRepository.UserFilter{
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
	return user.UserDtoFromEntity(targetEntityUser), nil
}

func (srv *Service) GetUserByEmail(ctx context.Context, email string, password string) (*user.User, error) {
	targetEntityUser, ex, err := srv.userRepo.GetUser(ctx, userRepository.UserFilter{
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
	return user.UserDtoFromEntity(targetEntityUser), nil
}

func (srv *Service) comparePassword(origin, existed string) bool {
	return generatePasswordHash(origin) == existed
}

func generatePasswordHash(password string) string {
	hash := sha512.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
