package note

import (
	"context"
	"wn/internal/domain/dto/request"
	req "wn/internal/domain/dto/request"
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

type srv interface {
	CreateNote(ctx context.Context, req req.NoteRequest, userId uuid.UUID) (uuid.UUID, error)
	UpdateNote(ctx context.Context, req req.NoteWithIdRequest, userId uuid.UUID) error
	DeleteNote(ctx context.Context, req req.NoteId, userId uuid.UUID) error
}

type Controller struct {
	lgr     applogger.Logger
	builder *response.Builder

	noteService srv
}

func NewController(logger applogger.Logger, builder *response.Builder, noteService srv) *Controller {

	return &Controller{
		lgr:     logger,
		builder: builder,

		noteService: noteService,
	}
}

func (h *Controller) Init(api, authApi *gin.RouterGroup) {
	notesAuth := authApi.Group("/notes")
	{
		notesAuth.POST("/create", h.createNote)
		notesAuth.POST("/update", h.updateNote)
		notesAuth.POST("/delete", h.deleteNote)
	}
}

// @Summary create_note
// @Description Создать заметку
// @Tags notes
// @Produce json
// @Param data body request.NoteRequest true "data"
// @Param X-Request-Id header string true "Request id identity"
// @Param Authorization header string true "auth token"
// @Success 200 {object} response.Response{data=resp.NoteId}
// @Failure 400 {object} response.Response{} "possible codes: invalid_token, invalid_authorization_header"
// @Failure 400 {object} response.Response{} "possible codes: bind_body, invalid_X-Request-Id"
// @Failure 422 {object} response.Response{} "possible codes: not_unique"
// @Router /wn/api/v1/notes/create [post]
func (h *Controller) createNote(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.NoteRequest
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

	noteId, err := h.noteService.CreateNote(ctx, req, userId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, resp.NoteId{
		Id: noteId,
	}))
}

// @Summary delete_note
// @Description Удалить заметку
// @Tags notes
// @Produce json
// @Param data body request.NoteId true "data"
// @Param X-Request-Id header string true "Request id identity"
// @Param Authorization header string true "auth token"
// @Success 200 {object} response.Response{}
// @Failure 400 {object} response.Response{} "possible codes: invalid_token, invalid_authorization_header"
// @Failure 400 {object} response.Response{} "possible codes: bind_body, invalid_X-Request-Id"
// @Failure 422 {object} response.Response{} "possible codes: note_not_found"
// @Router /wn/api/v1/notes/delete [post]
func (h *Controller) deleteNote(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.NoteId
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

	err = h.noteService.DeleteNote(ctx, req, userId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, nil))
}

// @Summary update_note
// @Description Обновить заметку
// @Tags notes
// @Produce json
// @Param data body request.NoteWithIdRequest true "data"
// @Param X-Request-Id header string true "Request id identity"
// @Param Authorization header string true "auth token"
// @Success 200 {object} response.Response{}
// @Failure 400 {object} response.Response{} "possible codes: invalid_token, invalid_authorization_header"
// @Failure 400 {object} response.Response{} "possible codes: bind_body, invalid_X-Request-Id"
// @Failure 422 {object} response.Response{} "possible codes: note_not_found"
// @Router /wn/api/v1/notes/update [post]
func (h *Controller) updateNote(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.NoteWithIdRequest
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

	err = h.noteService.UpdateNote(ctx, req, userId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, nil))
}
