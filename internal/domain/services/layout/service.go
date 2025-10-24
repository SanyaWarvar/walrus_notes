package layout

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

type layoutRepo interface {
	CreateLayout(ctx context.Context, item *entity.Layout) (uuid.UUID, error)
	DeleteLayoutById(ctx context.Context, layoutId, userId uuid.UUID) error
	GetAvailableLayouts(ctx context.Context, userId uuid.UUID) ([]entity.Layout, error)
}

type linksRepo interface {
	DeleteLinkNotes(ctx context.Context, layoutId, firstNoteId, secondNoteId uuid.UUID) error
	DeleteLinksFromLayout(ctx context.Context, layoutId uuid.UUID) error
	DeleteLinksWithNote(ctx context.Context, noteId uuid.UUID) error
	LinkNotes(ctx context.Context, layoutId, firstNoteId, secondNoteId uuid.UUID) error
	GetAllLinks(ctx context.Context, layoutId uuid.UUID, noteIds []uuid.UUID) ([]entity.Link, error)
	DeleteLayoutNote(ctx context.Context, layoutId uuid.UUID, userId uuid.UUID) error
}

type Service struct {
	tx     trx.TransactionManager
	logger applogger.Logger

	layoutRepo layoutRepo
	linksRepo  linksRepo
}

func NewService(
	tx trx.TransactionManager,
	logger applogger.Logger,
	layoutRepo layoutRepo,
	linksRepo linksRepo,
) *Service {
	return &Service{
		tx:         tx,
		logger:     logger,
		layoutRepo: layoutRepo,
		linksRepo:  linksRepo,
	}
}

func (srv *Service) CreateLayout(ctx context.Context, title string, ownerId uuid.UUID) (uuid.UUID, error) {
	item := entity.Layout{
		Id:         util.NewUUID(),
		Title:      title,
		OwnerId:    ownerId,
		HaveAccess: []uuid.UUID{ownerId},
	}
	return srv.layoutRepo.CreateLayout(ctx, &item)
}

func (srv *Service) DeleteLayoutById(ctx context.Context, layoutId, ownerId uuid.UUID) error {
	return srv.tx.Transaction(ctx, func(ctx context.Context) error {
		err := srv.linksRepo.DeleteLayoutNote(ctx, layoutId, ownerId)
		if err != nil {
			return errors.Wrap(err, "srv.linksRepo.DeleteLayoutNote")
		}

		err = srv.linksRepo.DeleteLinksFromLayout(ctx, layoutId)
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
		})
	}

	return output, nil
}
