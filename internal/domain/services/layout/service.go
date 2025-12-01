package layout

import (
	"context"
	"wn/internal/domain/dto"
	"wn/internal/domain/dto/request"
	"wn/internal/entity"
	apperrors "wn/internal/errors"
	"wn/pkg/applogger"
	"wn/pkg/trx"
	"wn/pkg/util"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type layoutRepo interface {
	CreateLayout(ctx context.Context, item *entity.Layout) (uuid.UUID, error)
	DeleteLayoutById(ctx context.Context, layoutId, userId uuid.UUID) error
	GetAvailableLayouts(ctx context.Context, userId uuid.UUID) ([]entity.Layout, error)
	UpdateLayout(ctx context.Context, userId, layoutId uuid.UUID, color, title string) (int, error)
}

type linksRepo interface {
	DeleteLinksWithNote(ctx context.Context, noteId uuid.UUID) error
	LinkNotes(ctx context.Context, firstNoteId, secondNoteId uuid.UUID) error
	GetAllLinks(ctx context.Context, noteIds []uuid.UUID) ([]entity.Link, error)
}

type noteRepo interface {
	DeleteNotesByLayoutId(ctx context.Context, layoutId, userId uuid.UUID) error
}

type Service struct {
	tx     trx.TransactionManager
	logger applogger.Logger

	layoutRepo layoutRepo
	linksRepo  linksRepo
	noteRepo   noteRepo
}

func NewService(
	tx trx.TransactionManager,
	logger applogger.Logger,
	layoutRepo layoutRepo,
	linksRepo linksRepo,
	noteRepo noteRepo,
) *Service {
	return &Service{
		tx:         tx,
		logger:     logger,
		layoutRepo: layoutRepo,
		linksRepo:  linksRepo,
		noteRepo:   noteRepo,
	}
}

func (srv *Service) CreateLayout(ctx context.Context, title, color string, ownerId uuid.UUID, isMain bool) (uuid.UUID, error) {
	item := entity.Layout{
		Id:         util.NewUUID(),
		Title:      title,
		OwnerId:    ownerId,
		HaveAccess: []uuid.UUID{ownerId},
		IsMain:     isMain,
		Color:      color,
	}
	return srv.layoutRepo.CreateLayout(ctx, &item)
}

func (srv *Service) DeleteLayoutById(ctx context.Context, layoutId, ownerId uuid.UUID) error {
	return srv.tx.Transaction(ctx, func(ctx context.Context) error {

		err := srv.noteRepo.DeleteNotesByLayoutId(ctx, layoutId, ownerId)
		if err != nil {
			return errors.Wrap(err, "srv.linksRepo.DeleteLinksFromLayout")
		}

		err = srv.layoutRepo.DeleteLayoutById(ctx, layoutId, ownerId)
		if err != nil {
			return errors.Wrap(err, "srv.layoutRepo.DeleteLayoutById")
		}

		return nil
	})
}

func (srv *Service) GetAvailableLayouts(ctx context.Context, userId uuid.UUID) ([]dto.Layout, error) {
	entities, err := srv.layoutRepo.GetAvailableLayouts(ctx, userId)
	if err != nil {
		return nil, err
	}

	output := make([]dto.Layout, 0, len(entities))
	for _, item := range entities {
		output = append(output, dto.Layout{
			Id:      item.Id,
			OwnerId: item.OwnerId,
			Title:   item.Title,
			IsMain:  item.IsMain,
			Color:   item.Color,
		})
	}

	return output, nil
}

func (srv *Service) UpdateLayout(ctx context.Context, req request.UpdateLayout, userId uuid.UUID) error {
	updatedRows, err := srv.layoutRepo.UpdateLayout(ctx, userId, req.LayoutId, req.Color, req.Title)
	if err != nil {
		return errors.Wrap(err, "srv.layoutRepo.UpdateLayout")
	}
	if updatedRows < 1 {
		return apperrors.LayoutNotFound
	}
	return nil
}
