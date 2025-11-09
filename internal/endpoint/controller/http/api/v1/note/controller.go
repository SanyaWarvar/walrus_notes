package note

import (
	"context"
	"strconv"
	"wn/internal/domain/dto"
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
	CreateNote(ctx context.Context, req req.NoteRequest, userId uuid.UUID, mainLayoutId uuid.UUID) (uuid.UUID, error)
	UpdateNote(ctx context.Context, req req.NoteWithIdRequest, userId uuid.UUID) error
	DeleteNote(ctx context.Context, req req.NoteId, userId uuid.UUID, mainLayoutId uuid.UUID) error
	GetNotesFromLayout(ctx context.Context, req req.GetNotesFromLayoutRequest, userId uuid.UUID) ([]dto.Note, int, error)
	GetNotesWithPosition(ctx context.Context, userId uuid.UUID, req req.GetNotesFromLayoutWithoutPagRequest) ([]dto.Note, error)
	GetNotesWithoutPosition(ctx context.Context, userId uuid.UUID, req req.GetNotesFromLayoutWithoutPagRequest) ([]dto.Note, error)
	UpdateNotePosition(ctx context.Context, userId uuid.UUID, req req.UpdateNotePositionRequest) error
	SearchNotes(ctx context.Context, userId uuid.UUID, search string) ([]dto.Note, error)

	CreateLink(ctx context.Context, userId uuid.UUID, req req.LinkBetweenNotesRequest) error
	DeleteLink(ctx context.Context, userId uuid.UUID, req req.LinkBetweenNotesRequest) error
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
		notesAuth.GET("/search", h.searchNotes)
		layout := notesAuth.Group("/layout")
		{
			layout.GET("", h.getNotesFromLayout)
			layout.GET("/graph/posed", h.getNotesFromLayoutWithPosition)
			layout.GET("/graph/unposed", h.getNotesFromLayoutWithoutPosition)
			layout.POST("/graph/note", h.updateNotePosition)
		}

		links := layout.Group("/links")
		{
			links.POST("/create", h.createLinkBetweenNotes)
			links.POST("/delete", h.deleteLinkBetweenNotes)
		}
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

	mainLayoutId, err := uuid.Parse(ctx.Value("mainLayoutId").(string))
	if err != nil {
		_ = c.Error(apperrors.InvalidAuthorizationHeader)
		return
	}

	noteId, err := h.noteService.CreateNote(ctx, req, userId, mainLayoutId)
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

	mainLayoutId, err := uuid.Parse(ctx.Value("mainLayoutId").(string))
	if err != nil {
		_ = c.Error(apperrors.InvalidAuthorizationHeader)
		return
	}

	err = h.noteService.DeleteNote(ctx, req, userId, mainLayoutId)
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

// @Summary get_notes_from_layout
// @Description Получить все заметки из layout-а
// @Tags notes
// @Produce json
// @Param page query int true "page"
// @Param layoutId query string true "layoutId"
// @Param X-Request-Id header string true "Request id identity"
// @Param Authorization header string true "auth token"
// @Success 200 {object} response.Response{}
// @Failure 400 {object} response.Response{} "possible codes: invalid_token, invalid_authorization_header"
// @Failure 400 {object} response.Response{} "possible codes: bind_query, invalid_X-Request-Id"
// @Failure 422 {object} response.Response{} "possible codes: note_not_found"
// @Router /wn/api/v1/notes/layout [get]
func (h *Controller) getNotesFromLayout(c *gin.Context) {
	ctx := c.Request.Context()
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.BindQueryError))
		return
	}

	layoutId, err := uuid.Parse(c.Query("layoutId"))
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.BindQueryError))
		return
	}

	req := request.GetNotesFromLayoutRequest{
		Page:     page,
		LayoutId: layoutId,
	}

	userId, err := util.GetUserId(ctx)
	if err != nil {
		_ = c.Error(apperrors.InvalidAuthorizationHeader)
		return
	}

	notes, count, err := h.noteService.GetNotesFromLayout(ctx, req, userId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(
		200,
		h.builder.BuildSuccessPaginationResponse(
			ctx,
			req.Page,
			constants.PageSize,
			count/constants.PageSize,
			notes,
		))
}

// @Summary get_unposed_notes
// @Description Получить заметки, которые не имеют позиции в графе
// @Tags graph
// @Produce json
// @Param layoutId query string true "layoutId"
// @Param X-Request-Id header string true "Request id identity"
// @Param Authorization header string true "auth token"
// @Success 200 {object} response.Response{}
// @Failure 400 {object} response.Response{} "possible codes: invalid_token, invalid_authorization_header"
// @Failure 400 {object} response.Response{} "possible codes: bind_body, invalid_X-Request-Id"
// @Router /wn/api/v1/notes/layout/graph/unposed [get]
func (h *Controller) getNotesFromLayoutWithoutPosition(c *gin.Context) {
	ctx := c.Request.Context()
	layoutId, err := uuid.Parse(c.Query("layoutId"))
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.BindQueryError))
		return
	}
	var req request.GetNotesFromLayoutWithoutPagRequest
	req.LayoutId = layoutId

	userId, err := util.GetUserId(ctx)
	if err != nil {
		_ = c.Error(apperrors.InvalidAuthorizationHeader)
		return
	}

	notes, err := h.noteService.GetNotesWithoutPosition(ctx, userId, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, notes))
}

// @Summary get_posed_notes
// @Description Получить заметки, которые имеют позиции в графе
// @Tags graph
// @Produce json
// @Param layoutId query string true "layoutId"
// @Param X-Request-Id header string true "Request id identity"
// @Param Authorization header string true "auth token"
// @Success 200 {object} response.Response{}
// @Failure 400 {object} response.Response{} "possible codes: invalid_token, invalid_authorization_header"
// @Failure 400 {object} response.Response{} "possible codes: bind_body, invalid_X-Request-Id"
// @Router /wn/api/v1/notes/layout/graph/posed [get]
func (h *Controller) getNotesFromLayoutWithPosition(c *gin.Context) {
	ctx := c.Request.Context()

	layoutId, err := uuid.Parse(c.Query("layoutId"))
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.BindQueryError))
		return
	}
	var req request.GetNotesFromLayoutWithoutPagRequest
	req.LayoutId = layoutId

	userId, err := util.GetUserId(ctx)
	if err != nil {
		_ = c.Error(apperrors.InvalidAuthorizationHeader)
		return
	}

	notes, err := h.noteService.GetNotesWithPosition(ctx, userId, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, notes))
}

// @Summary update_note_position
// @Description Обновить позицию заметки в графе
// @Tags graph
// @Produce json
// @Param data body request.UpdateNotePositionRequest true "data"
// @Param X-Request-Id header string true "Request id identity"
// @Param Authorization header string true "auth token"
// @Success 200 {object} response.Response{}
// @Failure 400 {object} response.Response{} "possible codes: invalid_token, invalid_authorization_header"
// @Failure 400 {object} response.Response{} "possible codes: bind_body, invalid_X-Request-Id"
// @Router /wn/api/v1/notes/layout/graph/note [post]
func (h *Controller) updateNotePosition(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.UpdateNotePositionRequest
	err := c.BindJSON(&req)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.BindBodyError))
		return
	}

	userId, err := util.GetUserId(ctx)
	if err != nil {
		_ = c.Error(apperrors.InvalidAuthorizationHeader)
		return
	}

	err = h.noteService.UpdateNotePosition(ctx, userId, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, nil))
}

// @Summary create_link_between_notes
// @Description Связать заметки
// @Tags links
// @Produce json
// @Param data body request.LinkBetweenNotesRequest true "data"
// @Param X-Request-Id header string true "Request id identity"
// @Param Authorization header string true "auth token"
// @Success 200 {object} response.Response{}
// @Failure 400 {object} response.Response{} "possible codes: invalid_token, invalid_authorization_header"
// @Failure 400 {object} response.Response{} "possible codes: bind_body, invalid_X-Request-Id"
// @Router /wn/api/v1/notes/layout/links/create [post]
func (h *Controller) createLinkBetweenNotes(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.LinkBetweenNotesRequest
	err := c.BindJSON(&req)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.BindBodyError))
		return
	}

	userId, err := util.GetUserId(ctx)
	if err != nil {
		_ = c.Error(apperrors.InvalidAuthorizationHeader)
		return
	}

	err = h.noteService.CreateLink(ctx, userId, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, nil))
}

// @Summary delete_link_between_notes
// @Description Обновить позицию заметки в графе
// @Tags links
// @Produce json
// @Param data body request.LinkBetweenNotesRequest true "data"
// @Param X-Request-Id header string true "Request id identity"
// @Param Authorization header string true "auth token"
// @Success 200 {object} response.Response{}
// @Failure 400 {object} response.Response{} "possible codes: invalid_token, invalid_authorization_header"
// @Failure 400 {object} response.Response{} "possible codes: bind_body, invalid_X-Request-Id"
// @Router /wn/api/v1/notes/layout/links/delete [post]
func (h *Controller) deleteLinkBetweenNotes(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.LinkBetweenNotesRequest
	err := c.BindJSON(&req)
	if err != nil {
		_ = c.Error(apperror.NewBadRequestError(err.Error(), constants.BindBodyError))
		return
	}

	userId, err := util.GetUserId(ctx)
	if err != nil {
		_ = c.Error(apperrors.InvalidAuthorizationHeader)
		return
	}

	err = h.noteService.DeleteLink(ctx, userId, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, nil))
}

// @Summary search_notes
// @Description Найти заметку
// @Tags notes
// @Produce json
// @Param search query string true "search"
// @Param X-Request-Id header string true "Request id identity"
// @Param Authorization header string true "auth token"
// @Success 200 {object} response.Response{}
// @Failure 400 {object} response.Response{} "possible codes: invalid_token, invalid_authorization_header"
// @Failure 400 {object} response.Response{} "possible codes: bind_body, invalid_X-Request-Id"
// @Router /wn/api/v1/notes/search [post]
func (h *Controller) searchNotes(c *gin.Context) {
	ctx := c.Request.Context()

	userId, err := util.GetUserId(ctx)
	if err != nil {
		_ = c.Error(apperrors.InvalidAuthorizationHeader)
		return
	}

	search := c.Query("search")

	notes, err := h.noteService.SearchNotes(ctx, userId, search)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.AbortWithStatusJSON(h.builder.BuildSuccessResponseBody(ctx, notes))
}
