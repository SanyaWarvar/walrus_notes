package note

import (
	"context"
	"wn/internal/entity"
	"wn/pkg/applogger"
	"wn/pkg/trx"
	"wn/pkg/util"

	"github.com/google/uuid"
)

type noteRepo interface {
	DeleteNoteById(ctx context.Context, noteId, userId uuid.UUID) error
	CreateNote(ctx context.Context, item *entity.Note) (uuid.UUID, error)
	UpdateNote(ctx context.Context, newItem *entity.Note) error
}

type Service struct {
	tx     trx.TransactionManager
	logger applogger.Logger

	noteRepo noteRepo
}

func NewService(
	tx trx.TransactionManager,
	logger applogger.Logger,
	noteRepo noteRepo,
) *Service {
	return &Service{
		tx:       tx,
		logger:   logger,
		noteRepo: noteRepo,
	}
}

func (srv *Service) DeleteNoteById(ctx context.Context, noteId, userId uuid.UUID) error {
	return srv.noteRepo.DeleteNoteById(ctx, noteId, userId)
}

func (srv *Service) CreateNote(ctx context.Context, title, payload string, ownerId uuid.UUID) (uuid.UUID, error) {
	n := entity.Note{
		Id:         util.NewUUID(),
		Title:      title,
		Payload:    payload,
		CreatedAt:  util.GetCurrentUTCTime(),
		OwnerId:    ownerId,
		HaveAccess: []uuid.UUID{ownerId},
	}
	return srv.noteRepo.CreateNote(ctx, &n)
}

func (srv *Service) UpdateNote(ctx context.Context, title, payload string, noteId, ownerId uuid.UUID) error {
	n := entity.Note{
		Id:      noteId,
		Title:   title,
		Payload: payload,
		OwnerId: ownerId,
	}
	return srv.noteRepo.UpdateNote(ctx, &n)
}
