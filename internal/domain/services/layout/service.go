package layout

import (
	"context"
	"wn/internal/domain/dto"
	"wn/internal/entity"
	"wn/pkg/applogger"
	"wn/pkg/trx"
	"wn/pkg/util"

	"github.com/google/uuid"
)

type layoutRepo interface {
	CreateLayout(ctx context.Context, item *entity.Layout) (uuid.UUID, error)
	DeleteLayoutById(ctx context.Context, layoutId, userId uuid.UUID) error
	GetAvailableLayouts(ctx context.Context, userId uuid.UUID) ([]entity.Layout, error)
}

type Service struct {
	tx     trx.TransactionManager
	logger applogger.Logger

	layoutRepo layoutRepo
}

func NewService(
	tx trx.TransactionManager,
	logger applogger.Logger,
	layoutRepo layoutRepo,
) *Service {
	return &Service{
		tx:         tx,
		logger:     logger,
		layoutRepo: layoutRepo,
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
	return srv.layoutRepo.DeleteLayoutById(ctx, layoutId, ownerId)
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
