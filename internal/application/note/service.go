package note

import (
	"context"
	"wn/internal/domain/dto"
	req "wn/internal/domain/dto/request"
	"wn/pkg/applogger"
	"wn/pkg/trx"

	"github.com/google/uuid"
)

type noteService interface {
	DeleteNoteById(ctx context.Context, noteId, userId uuid.UUID) error
	CreateNote(ctx context.Context, title, payload string, ownerId, layoutId uuid.UUID) (uuid.UUID, error)
	UpdateNote(ctx context.Context, title, payload string, noteId, ownerId uuid.UUID) error
	GetNotesWithPagination(ctx context.Context, page int, layoutId, userId uuid.UUID) ([]dto.Note, int, error)
	GetNotesWithPosition(ctx context.Context, layoutId, userId uuid.UUID) ([]dto.Note, error)
	GetNotesWithoutPosition(ctx context.Context, layoutId, userId uuid.UUID) ([]dto.Note, error)
	UpdateNotePosition(ctx context.Context, layoutId, noteId uuid.UUID, xPos, yPos *float64) error
	CreateLink(ctx context.Context, layoutId, noteId1, noteId2 uuid.UUID) error
	DeleteLink(ctx context.Context, layoutId, noteId1, noteId2 uuid.UUID) error
}

type Service struct {
	tx     trx.TransactionManager
	logger applogger.Logger

	noteService noteService
}

func NewService(
	tx trx.TransactionManager,
	logger applogger.Logger,
	noteService noteService,
) *Service {
	return &Service{
		tx:          tx,
		logger:      logger,
		noteService: noteService,
	}
}

func (srv *Service) CreateNote(ctx context.Context, req req.NoteRequest, userId uuid.UUID) (uuid.UUID, error) {
	return srv.noteService.CreateNote(ctx, req.Title, req.Payload, userId, req.LayoutId)
}

func (srv *Service) UpdateNote(ctx context.Context, req req.NoteWithIdRequest, userId uuid.UUID) error {
	return srv.noteService.UpdateNote(ctx, req.Title, req.Payload, req.NoteId, userId)
}

func (srv *Service) DeleteNote(ctx context.Context, req req.NoteId, userId uuid.UUID) error {
	return srv.noteService.DeleteNoteById(ctx, req.NoteId, userId)
}

func (srv *Service) GetNotesFromLayout(ctx context.Context, req req.GetNotesFromLayoutRequest, userId uuid.UUID) ([]dto.Note, int, error) {
	return srv.noteService.GetNotesWithPagination(ctx, req.Page, req.LayoutId, userId)
}

func (srv *Service) GetNotesWithPosition(ctx context.Context, userId uuid.UUID, req req.GetNotesFromLayoutWithoutPagRequest) ([]dto.Note, error) {
	return srv.noteService.GetNotesWithPosition(ctx, req.LayoutId, userId)
}

func (srv *Service) GetNotesWithoutPosition(ctx context.Context, userId uuid.UUID, req req.GetNotesFromLayoutWithoutPagRequest) ([]dto.Note, error) {
	return srv.noteService.GetNotesWithoutPosition(ctx, req.LayoutId, userId)
}

func (srv *Service) UpdateNotePosition(ctx context.Context, userId uuid.UUID, req req.UpdateNotePositionRequest) error {
	return srv.noteService.UpdateNotePosition(ctx, req.LayoutId, req.NoteId, req.XPos, req.YPos)
}

func (srv *Service) CreateLink(ctx context.Context, userId uuid.UUID, req req.LinkBetweenNotesRequest) error {
	return srv.noteService.CreateLink(ctx, req.LayoutId, req.FirstNoteId, req.SecondNoteId)
}

func (srv *Service) DeleteLink(ctx context.Context, userId uuid.UUID, req req.LinkBetweenNotesRequest) error {
	return srv.noteService.DeleteLink(ctx, req.LayoutId, req.FirstNoteId, req.SecondNoteId)
}
