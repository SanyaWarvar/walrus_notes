package layout

import (
	"context"
	"wn/internal/domain/dto"
	"wn/internal/domain/dto/request"
	resp "wn/internal/domain/dto/response"
	apperrors "wn/internal/errors"
	"wn/pkg/apperror"
	"wn/pkg/applogger"
	"wn/pkg/constants"
	"wn/pkg/response"
	"wn/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type layoutService interface {
	CreateLayout(ctx context.Context, req request.NewLayoutRequest, userId uuid.UUID) (uuid.UUID, error)
	GetLayoutsByUserId(ctx context.Context, userId uuid.UUID) ([]dto.Layout, error)
}

type Controller struct {
	lgr     applogger.Logger
	builder *response.Builder

	layoutService layoutService
}

func NewController(logger applogger.Logger, builder *response.Builder, layoutService layoutService) *Controller {

	return &Controller{
		lgr:     logger,
		builder: builder,

		layoutService: layoutService,
	}
}

func (h *Controller) Init(api, authApi *gin.RouterGroup) {
	notesAuth := authApi.Group("/layout")
	{
		notesAuth.POST("/create", h.createLayout)
		notesAuth.GET("/my", h.getMyLayouts)
	}
}

// @Summary create_layout
// @Description Создать новый layout
// @Tags layouts
// @Produce json
// @Param data body request.NewLayoutRequest true "data"
// @Param X-Request-Id header string true "Request id identity"
// @Param Authorization header string true "auth token"
// @Success 200 {object} response.Response{data=resp.NoteId}
// @Failure 400 {object} response.Response{} "possible codes: invalid_token, invalid_authorization_header"
// @Failure 400 {object} response.Response{} "possible codes: bind_body, invalid_X-Request-Id"
// @Router /wn/api/v1/layout/create [post]
func (h *Controller) createLayout(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.NewLayoutRequest
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

	noteId, err := h.layoutService.CreateLayout(ctx, req, userId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, resp.NoteId{
		Id: noteId,
	}))
}

// @Summary get_my_layouts
// @Description Получить все layout-ы, к которым имеет доступ пользователь
// @Tags layouts
// @Produce json
// @Param data body request.NoteId true "data"
// @Param X-Request-Id header string true "Request id identity"
// @Param Authorization header string true "auth token"
// @Success 200 {object} response.Response{}
// @Failure 400 {object} response.Response{} "possible codes: invalid_token, invalid_authorization_header"
// @Failure 400 {object} response.Response{} "possible codes: bind_body, invalid_X-Request-Id"
// @Router /wn/api/v1/layout/my [get]
func (h *Controller) getMyLayouts(c *gin.Context) {
	ctx := c.Request.Context()

	userId, err := util.GetUserId(ctx)
	if err != nil {
		_ = c.Error(apperrors.InvalidAuthorizationHeader)
		return
	}

	layouts, err := h.layoutService.GetLayoutsByUserId(ctx, userId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, layouts))
}
