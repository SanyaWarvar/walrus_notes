package auth

import (
	"context"
	"wn/internal/domain/dto/request"
	resp "wn/internal/domain/dto/response"
	"wn/internal/domain/enum"
	"wn/internal/domain/services/token"
	"wn/pkg/apperror"
	"wn/pkg/applogger"
	"wn/pkg/constants"
	"wn/pkg/response"

	"github.com/gin-gonic/gin"
)

type userService interface {
	RegisterUser(ctx context.Context, credentials request.RegisterCredentials) (*resp.RegisterResponse, error)
}

type authService interface {
	SendConfirmationCode(ctx context.Context, req request.LoginRequest, action enum.EmailCodeAction) (*resp.SendCodeResponse, error)
	ConfirmCode(ctx context.Context, req request.ConfimationCodeRequest) error
	Login(ctx context.Context, req request.LoginRequest) (*token.UserTokens, error)
	RefreshTokens(ctx context.Context, req token.UserTokens) (*token.UserTokens, error)
}

type Controller struct {
	lgr     applogger.Logger
	builder *response.Builder

	userService userService
	authService authService
}

func NewController(logger applogger.Logger, builder *response.Builder, userService userService, authService authService) *Controller {
	return &Controller{
		lgr:     logger,
		builder: builder,

		userService: userService,
		authService: authService,
	}
}

func (h *Controller) Init(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", h.register)
		auth.POST("/login", h.login)
		auth.POST("/code", h.sendCode)
		auth.POST("/confirm", h.confirmCode)
		auth.POST("/refresh", h.refreshTokens)
		auth.POST("/forgot", h.forgotPassword)
	}
}

// @Summary register_user
// @Description register new user
// @Tags auth
// @Produce json
// @Param data body request.RegisterCredentials true "data"
// @Param X-Request-Id header string true "Request id identity"
// @Success 200 {object} response.Response{data=resp.RegisterResponse}
// @Failure 400 {object} response.Response{} "possible codes: bind_body, invalid_X-Request-Id"
// @Failure 422 {object} response.Response{} "possible codes: not_unique"
// @Router /wn/api/v1/auth/register [post]
func (h *Controller) register(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.RegisterCredentials
	err := c.BindJSON(&req)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.BindBodyError))
		return
	}
	userId, err := h.userService.RegisterUser(ctx, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, userId))
}

// @Summary send_confirm_code
// @Description register new user
// @Tags auth
// @Produce json
// @Param data body request.LoginRequest true "data"
// @Param X-Request-Id header string true "Request id identity"
// @Success 200 {object} response.Response{data=resp.SendCodeResponse}
// @Failure 400 {object} response.Response{} "possible codes: bind_body, invalid_X-Request-Id"
// @Failure 400 {object} response.Response{} "possible codes: incorrect_password"
// @Failure 422 {object} response.Response{} "possible codes: user_not_found, confirm_code_already_send"
// @Router /wn/api/v1/auth/code [post]
func (h *Controller) sendCode(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.LoginRequest
	err := c.BindJSON(&req)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.BindBodyError))
		return
	}
	if req.Password == "" {
		_ = c.Error(apperror.NewBadRequestError("password cant be empty", constants.BindBodyError))
		return
	}
	resp, err := h.authService.SendConfirmationCode(ctx, req, enum.ConfirmCode)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, resp))
}

// @Summary confirm_code
// @Description Подтверждение кода для подтверждения почты, либо сброса пароля. Если сброс пароля, то newPassword обязательное поле.
// @Tags auth
// @Produce json
// @Param data body request.ConfimationCodeRequest true "data"
// @Param X-Request-Id header string true "Request id identity"
// @Success 200 {object} response.Response{}
// @Failure 400 {object} response.Response{} "possible codes: bind_body, invalid_X-Request-Id"
// @Failure 422 {object} response.Response{} "possible codes: user_not_found, confirm_code_incorrect, confirm_code_not_exist, no_new_password"
// @Router /wn/api/v1/auth/confirm [post]
func (h *Controller) confirmCode(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.ConfimationCodeRequest
	err := c.BindJSON(&req)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.BindBodyError))
		return
	}
	err = h.authService.ConfirmCode(ctx, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, nil))
}

// @Summary login
// @Description Получение access,refresh токенов по почте и паролю
// @Tags auth
// @Produce json
// @Param data body request.LoginRequest true "data"
// @Param X-Request-Id header string true "Request id identity"
// @Success 200 {object} response.Response{data=token.UserTokens}
// @Failure 400 {object} response.Response{} "possible codes: bind_body, invalid_X-Request-Id"
// @Failure 401 {object} response.Response{} "possible codes: incorrect_password"
// @Failure 422 {object} response.Response{} "possible codes: user_not_found "
// @Router /wn/api/v1/auth/login [post]
func (h *Controller) login(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.LoginRequest
	err := c.BindJSON(&req)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.BindBodyError))
		return
	}
	tokens, err := h.authService.Login(ctx, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, tokens))
}

// @Summary refresh_tokens
// @Description Получение access,refresh токенов по access, refresh токенам
// @Tags auth
// @Produce json
// @Param data body token.UserTokens true "data"
// @Param X-Request-Id header string true "Request id identity"
// @Success 200 {object} response.Response{data=token.UserTokens}
// @Failure 400 {object} response.Response{} "possible codes: bind_body, invalid_X-Request-Id"
// @Failure 422 {object} response.Response{} "possible codes: bad_refresh_token, bad_access_token, bad_token_claims, token_dont_exist, tokens_dont_match"
// @Router /wn/api/v1/auth/refresh [post]
func (h *Controller) refreshTokens(c *gin.Context) {
	ctx := c.Request.Context()
	var req token.UserTokens
	err := c.BindJSON(&req)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.BindBodyError))
		return
	}
	tokens, err := h.authService.RefreshTokens(ctx, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, tokens))
}

// @Summary forgot_password
// @Description Сброс пароля
// @Tags auth
// @Produce json
// @Param data body request.LoginRequest true "data"
// @Param X-Request-Id header string true "Request id identity"
// @Success 200 {object} response.Response{data=resp.SendCodeResponse}
// @Failure 400 {object} response.Response{} "possible codes: bind_body, invalid_X-Request-Id"
// @Failure 422 {object} response.Response{} "possible codes: user_not_found, confirm_code_already_send"
// @Router /wn/api/v1/auth/forgot [post]
func (h *Controller) forgotPassword(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.LoginRequest
	err := c.BindJSON(&req)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.BindBodyError))
		return
	}
	resp, err := h.authService.SendConfirmationCode(ctx, req, enum.ForgotPassword)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, resp))
}
