package file

import (
	"context"
	"wn/internal/domain/dto"
	"wn/internal/domain/dto/request"
	apperrors "wn/internal/errors"
	"wn/pkg/apperror"
	"wn/pkg/applogger"
	"wn/pkg/constants"
	"wn/pkg/response"
	"wn/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type fileService interface {
	UploadFile(ctx context.Context, userId uuid.UUID, req request.UploadFileRequest, host string) (*dto.UploadFileResponse, error)
}

type Controller struct {
	lgr     applogger.Logger
	builder *response.Builder

	fileService fileService
}

func NewController(logger applogger.Logger, builder *response.Builder, fileService fileService) *Controller {
	return &Controller{
		lgr:     logger,
		builder: builder,

		fileService: fileService,
	}
}

func (h *Controller) Init(api, authApi *gin.RouterGroup) {
	_ = api.Group("/file")
	fileAuth := authApi.Group("/file")
	{
		fileAuth.POST("/upload", h.uploadFile)
	}
}

// todo add validation file

// @Summary upload_file
// @Description загрузить файл
// @Tags file
// @Produce json
// @Param data body request.UploadFileRequest true "data"
// @Param X-Request-Id header string true "Request id identity"
// @Param Authorization header string true "auth token"
// @Success 200 {object} response.Response{data=dto.UploadFileResponse}
// @Failure 400 {object} response.Response{} "possible codes: invalid_token, invalid_authorization_header"
// @Failure 400 {object} response.Response{} "possible codes: bind_body, invalid_X-Request-Id"
// @Router /wn/api/v1/file/upload [post]
func (h *Controller) uploadFile(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.UploadFileRequest
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

	picUrl, err := h.fileService.UploadFile(ctx, userId, req, c.Request.Host)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, picUrl))
}
