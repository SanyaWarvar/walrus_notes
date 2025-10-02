package note

import (
	"context"
	req "wn/internal/domain/dto/request"
	"wn/pkg/applogger"
	"wn/pkg/trx"

	"github.com/google/uuid"
)

type noteService interface {
	DeleteNoteById(ctx context.Context, noteId, userId uuid.UUID) error
	CreateNote(ctx context.Context, title, payload string, ownerId uuid.UUID) (uuid.UUID, error)
	UpdateNote(ctx context.Context, title, payload string, noteId, ownerId uuid.UUID) error
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
	return srv.noteService.CreateNote(ctx, req.Title, req.Payload, userId)
}

func (srv *Service) UpdateNote(ctx context.Context, req req.NoteWithIdRequest, userId uuid.UUID) error {
	return srv.noteService.UpdateNote(ctx, req.Title, req.Payload, req.NoteId, userId)
}

func (srv *Service) DeleteNote(ctx context.Context, req req.NoteId, userId uuid.UUID) error {
	return srv.noteService.DeleteNoteById(ctx, req.NoteId, userId)
}
