package auth

import (
	"context"
	"time"
	"wn/internal/domain/dto"
	"wn/internal/domain/enum"
	"wn/internal/domain/services/token"
	apperrors "wn/internal/errors"
	"wn/pkg/applogger"
	"wn/pkg/trx"

	"github.com/google/uuid"
)

var codeDelay time.Duration = time.Duration(time.Minute * 1)

type userService interface {
	GetUserByEmail(ctx context.Context, email string, password string) (*dto.User, error)
	UpdateUser(ctx context.Context, userId uuid.UUID, filter *dto.UserUpdateParams) error
}

type tokenService interface {
	GenerateUserTokens(ctx context.Context, userId, mainLayoutId uuid.UUID, role string) (*token.UserTokens, error)
	ParseToken(token string, withExpCheck bool) (*token.CustomClaims, error)
	RefreshTokens(ctx context.Context, access, refresh string) (*token.UserTokens, error)
}

type smtpService interface {
	SendConfirmEmailCode(ctx context.Context, email string, action enum.EmailCodeAction) error
	ConfirmCode(ctx context.Context, email string, code string) (*dto.ConfirmationCode, error)
}

type layoutService interface {
	GetAvailableLayouts(ctx context.Context, userId uuid.UUID) ([]dto.Layout, error)
}

type Service struct {
	tx     trx.TransactionManager
	logger applogger.Logger

	userService   userService
	smtpService   smtpService
	tokenService  tokenService
	layoutService layoutService
}

func NewService(
	tx trx.TransactionManager,
	logger applogger.Logger,
	userService userService,
	smtpService smtpService,
	tokenService tokenService,
	layoutService layoutService,
) *Service {
	return &Service{
		tx:            tx,
		logger:        logger,
		userService:   userService,
		smtpService:   smtpService,
		tokenService:  tokenService,
		layoutService: layoutService,
	}
}

// todo add check is confirmed email

func (srv *Service) SendConfirmationCode(ctx context.Context, req dto.LoginRequest, action enum.EmailCodeAction) (*dto.SendCodeResponse, error) {
	_, err := srv.userService.GetUserByEmail(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &dto.SendCodeResponse{NextCodeDelay: codeDelay},
		srv.smtpService.SendConfirmEmailCode(ctx, req.Email, action)
}

func (srv *Service) ConfirmCode(ctx context.Context, req dto.ConfirmationCodeRequest) error {
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
		return srv.userService.UpdateUser(ctx, u.Id, &dto.UserUpdateParams{
			ConfirmedEmail: &t,
		})
	case enum.ForgotPassword:
		if req.NewPassword == "" {
			return apperrors.NoNewPassword
		}
		return srv.userService.UpdateUser(ctx, u.Id, &dto.UserUpdateParams{
			Password: &req.NewPassword,
		})
	default:
		return nil
	}
}

func (srv *Service) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	u, err := srv.userService.GetUserByEmail(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	layouts, err := srv.layoutService.GetAvailableLayouts(ctx, u.Id)
	var layoutId uuid.UUID
	for _, layout := range layouts {
		if layout.OwnerId == u.Id && layout.IsMain {
			layoutId = layout.Id
			break
		}
	}

	tokens, err := srv.tokenService.GenerateUserTokens(ctx, u.Id, layoutId, u.Role)
	if err != nil {
		return nil, err
	}
	return &dto.LoginResponse{
		UserId:  u.Id,
		Access:  tokens.Access,
		Refresh: tokens.Refresh,
	}, nil
}

func (srv *Service) RefreshTokens(ctx context.Context, req token.UserTokens) (*token.UserTokens, error) {
	return srv.tokenService.RefreshTokens(ctx, req.Access, req.Refresh)
}
