package auth

import (
	"context"
	"wn/internal/domain/dto/auth"
	"wn/internal/domain/dto/request"
	dto "wn/internal/domain/dto/response"
	"wn/internal/domain/dto/user"
	"wn/internal/domain/enum"
	"wn/internal/domain/services/token"
	apperrors "wn/internal/errors"
	userRepository "wn/internal/infrastructure/repository/user"
	"wn/pkg/applogger"
	"wn/pkg/trx"
	"time"

	"github.com/google/uuid"
)

var codeDelay time.Duration = time.Duration(time.Minute * 1)

type userService interface {
	GetUserByEmail(ctx context.Context, email string, password string) (*user.User, error)
	UpdateUser(ctx context.Context, userId uuid.UUID, filter *userRepository.UserUpdateParams) error
}

type tokenService interface {
	GenerateUserTokens(ctx context.Context, id uuid.UUID, role string) (*token.UserTokens, error)
	ParseToken(token string) (*token.CustomClaims, error)
	RefreshTokens(ctx context.Context, access, refresh string) (*token.UserTokens, error)
}

type smtpService interface {
	SendConfirmEmailCode(ctx context.Context, email string, action enum.EmailCodeAction) error
	ConfirmCode(ctx context.Context, email string, code string) (*auth.ConfirmationCode, error)
}

type Service struct {
	tx     trx.TransactionManager
	logger applogger.Logger

	userService  userService
	smtpService  smtpService
	tokenService tokenService
}

func NewService(
	tx trx.TransactionManager,
	logger applogger.Logger,
	userService userService,
	smtpService smtpService,
	tokenService tokenService,
) *Service {
	return &Service{
		tx:           tx,
		logger:       logger,
		userService:  userService,
		smtpService:  smtpService,
		tokenService: tokenService,
	}
}

// todo add check is confirmed email

func (srv *Service) SendConfirmationCode(ctx context.Context, req request.LoginRequest, action enum.EmailCodeAction) (*dto.SendCodeResponse, error) {
	_, err := srv.userService.GetUserByEmail(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &dto.SendCodeResponse{NextCodeDelay: codeDelay},
		srv.smtpService.SendConfirmEmailCode(ctx, req.Email, action)
}

func (srv *Service) ConfirmCode(ctx context.Context, req request.ConfimationCodeRequest) error {
	u, err := srv.userService.GetUserByEmail(ctx, req.Email, "")
	if err != nil {
		return err
	}
	code, err := srv.smtpService.ConfirmCode(ctx, req.Email, req.Code)
	if err != nil {
		return err
	}
	t := true
	switch code.Action {
	case enum.ConfirmCode:
		return srv.userService.UpdateUser(ctx, u.Id, &userRepository.UserUpdateParams{
			ConfirmedEmail: &t,
		})
	case enum.ForgotPassword:
		if req.NewPassword == "" {
			return apperrors.NoNewPassword
		}
		return srv.userService.UpdateUser(ctx, u.Id, &userRepository.UserUpdateParams{
			Password: &req.NewPassword,
		})
	default:
		return nil
	}
}

func (srv *Service) Login(ctx context.Context, req request.LoginRequest) (*token.UserTokens, error) {
	u, err := srv.userService.GetUserByEmail(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return srv.tokenService.GenerateUserTokens(ctx, u.Id, u.Role)
}

func (srv *Service) RefreshTokens(ctx context.Context, req token.UserTokens) (*token.UserTokens, error) {
	return srv.tokenService.RefreshTokens(ctx, req.Access, req.Refresh)
}
