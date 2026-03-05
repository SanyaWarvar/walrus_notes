package permissions

import (
	"context"
	"wn/internal/domain/dto"
	apperrors "wn/internal/errors"
	"wn/pkg/apperror"
	"wn/pkg/applogger"
	"wn/pkg/constants"
	"wn/pkg/response"
	"wn/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type permissionsService interface {
	GeneratePermissionsLink(ctx context.Context, userId uuid.UUID, req *dto.GeneratePermissionLinkRequest) (*dto.GeneratePermissionsLinkResponse, error)
	ApplyPermissionsLink(ctx context.Context, userId uuid.UUID, req *dto.ApplyPermissionsRequest) error
	GetPermissionsDashboard(ctx context.Context, userId uuid.UUID) (*dto.PermissionsDashbord, error)
}

type Controller struct {
	lgr     applogger.Logger
	builder *response.Builder

	permissionsService permissionsService
}

func NewController(logger applogger.Logger, builder *response.Builder, permissionsService permissionsService) *Controller {
	return &Controller{
		lgr:     logger,
		builder: builder,

		permissionsService: permissionsService,
	}
}

func (h *Controller) Init(api, authApi *gin.RouterGroup) {
	permissionsAuth := authApi.Group("/permissions")
	{
		permissionsAuth.POST("/links/generate", h.generateLink)
		permissionsAuth.POST("/links/apply", h.applyLink)
		permissionsAuth.GET("/dashboard", h.getDashboard)
	}
}

// @Summary generateLink
// @Description Сгенерировать id для выдачи пермишенов
// @Tags permissions
// @Produce json
// @Param data body dto.GeneratePermissionLinkRequest true "data"
// @Param X-Request-Id header string true "Request id identity"
// @Param Authorization header string true "auth token"
// @Success 200 {object} response.Response{data=dto.GeneratePermissionsLinkResponse}
// @Failure 400 {object} response.Response{} "possible codes: invalid_token, invalid_authorization_header"
// @Failure 400 {object} response.Response{} "possible codes: bind_body, invalid_X-Request-Id"
// @Failure 422 {object} response.Response{} "possible codes: bad_kind, premissions_not_enough"
// @Router /wn/api/v1/links/generate [post]
func (h *Controller) generateLink(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.GeneratePermissionLinkRequest
	err := c.ShouldBind(&req)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.BindBodyError))
		return
	}
	userId, err := util.GetUserId(ctx)
	if err != nil {
		_ = c.Error(apperrors.InvalidAuthorizationHeader)
		return
	}

	resp, err := h.permissionsService.GeneratePermissionsLink(ctx, userId, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, resp))
}

// @Summary applyLink
// @Description Сгенерировать id для выдачи пермишенов
// @Tags permissions
// @Produce json
// @Param data body dto.ApplyPermissionsRequest true "data"
// @Param X-Request-Id header string true "Request id identity"
// @Param Authorization header string true "auth token"
// @Success 200 {object} response.Response{}
// @Failure 400 {object} response.Response{} "possible codes: invalid_token, invalid_authorization_header"
// @Failure 400 {object} response.Response{} "possible codes: bind_body, invalid_X-Request-Id"
// @Failure 422 {object} response.Response{} "possible codes: already_exists, cant_apply"
// @Router /wn/api/v1/links/apply [post]
func (h *Controller) applyLink(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.ApplyPermissionsRequest
	err := c.ShouldBind(&req)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.BindBodyError))
		return
	}
	userId, err := util.GetUserId(ctx)
	if err != nil {
		_ = c.Error(apperrors.InvalidAuthorizationHeader)
		return
	}

	err = h.permissionsService.ApplyPermissionsLink(ctx, userId, &req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, nil))
}

// @Summary getDashboard
// @Description Получить дашборд пермишенов
// @Tags permissions
// @Produce json
// @Param X-Request-Id header string true "Request id identity"
// @Param Authorization header string true "auth token"
// @Success 200 {object} response.Response{data=dto.PermissionsDashbord}
// @Failure 400 {object} response.Response{} "possible codes: invalid_token, invalid_authorization_header"
// @Failure 400 {object} response.Response{} "possible codes: invalid_X-Request-Id"
// @Router /wn/api/v1/dashboard [get]
func (h *Controller) getDashboard(c *gin.Context) {
	ctx := c.Request.Context()

	userId, err := util.GetUserId(ctx)
	if err != nil {
		_ = c.Error(apperrors.InvalidAuthorizationHeader)
		return
	}

	dashboard, err := h.permissionsService.GetPermissionsDashboard(ctx, userId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, dashboard))
}
