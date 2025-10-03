package note

import (
	"context"
	"wn/internal/domain/dto"
	"wn/internal/entity"
	"wn/pkg/applogger"
	"wn/pkg/trx"
	"wn/pkg/util"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type noteRepo interface {
	DeleteNoteById(ctx context.Context, noteId, userId uuid.UUID) error
	CreateNote(ctx context.Context, item *entity.Note) (uuid.UUID, error)
	UpdateNote(ctx context.Context, newItem *entity.Note) error
	GetNoteCountInLayout(ctx context.Context, layoutId uuid.UUID) (int, error)
	GetNotesByLayoutId(ctx context.Context, layoutId, userId uuid.UUID, offset, limit int) ([]entity.Note, error)
}

type layoutRepo interface {
	LinkNoteWithLayout(ctx context.Context, noteId, layoutId uuid.UUID) error
}

type Service struct {
	tx     trx.TransactionManager
	logger applogger.Logger

	noteRepo   noteRepo
	layoutRepo layoutRepo
}

func NewService(
	tx trx.TransactionManager,
	logger applogger.Logger,
	noteRepo noteRepo,
	layoutRepo layoutRepo,
) *Service {
	return &Service{
		tx:         tx,
		logger:     logger,
		noteRepo:   noteRepo,
		layoutRepo: layoutRepo,
	}
}

func (srv *Service) DeleteNoteById(ctx context.Context, noteId, userId uuid.UUID) error {
	return srv.noteRepo.DeleteNoteById(ctx, noteId, userId)
}

// todo add check access
func (srv *Service) CreateNote(ctx context.Context, title, payload string, ownerId, layoutId uuid.UUID) (uuid.UUID, error) {
	n := entity.Note{
		Id:         util.NewUUID(),
		Title:      title,
		Payload:    payload,
		CreatedAt:  util.GetCurrentUTCTime(),
		OwnerId:    ownerId,
		HaveAccess: []uuid.UUID{ownerId},
	}
	id, err := srv.noteRepo.CreateNote(ctx, &n)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "srv.noteRepo.CreateNote")
	}
	return id, srv.layoutRepo.LinkNoteWithLayout(ctx, id, layoutId)
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

func (srv *Service) GetNotesWithPagination(ctx context.Context, page int, layoutId, userId uuid.UUID) ([]dto.Note, int, error) {

	count, err := srv.noteRepo.GetNoteCountInLayout(ctx, layoutId)
	if err != nil {
		return nil, 0, errors.Wrap(err, "srv.noteRepo.GetNoteCountInLayout")
	}
	offset := util.CalculateOffset(page)
	limit := util.CalculateLimit()
	notes, err := srv.noteRepo.GetNotesByLayoutId(ctx, layoutId, userId, offset, limit)
	if err != nil {
		return nil, 0, errors.Wrap(err, "srv.noteRepo.GetNotesByLayoutId")
	}
	notesDto := dto.NotesFromEntities(notes)
	return notesDto, count, err
}
