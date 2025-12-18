package layout

import (
	"context"
	"wn/internal/domain/dto"
	"wn/internal/domain/dto/request"
	"wn/pkg/applogger"
	"wn/pkg/trx"

	"github.com/google/uuid"
)

type layoutService interface {
	CreateLayout(ctx context.Context, title, color string, ownerId uuid.UUID, isMain bool) (uuid.UUID, error)
	DeleteLayoutById(ctx context.Context, layoutId, ownerId uuid.UUID) error
	GetAvailableLayouts(ctx context.Context, userId uuid.UUID) ([]dto.Layout, error)
	ExportLayouts(ctx context.Context, userId uuid.UUID) (*dto.ExportInfo, error)
	UpdateLayout(ctx context.Context, req request.UpdateLayout, userId uuid.UUID) error
	ImportLayouts(ctx context.Context, userId uuid.UUID, info *dto.ExportInfo) error
}

type Service struct {
	tx     trx.TransactionManager
	logger applogger.Logger

	layoutService layoutService
}

func NewService(
	tx trx.TransactionManager,
	logger applogger.Logger,
	layoutService layoutService,
) *Service {
	return &Service{
		tx:            tx,
		logger:        logger,
		layoutService: layoutService,
	}
}

func (srv *Service) CreateLayout(ctx context.Context, req request.NewLayoutRequest, userId uuid.UUID) (uuid.UUID, error) {
	return srv.layoutService.CreateLayout(ctx, req.Title, req.Color, userId, false)
}

func (srv *Service) GetLayoutsByUserId(ctx context.Context, userId uuid.UUID) ([]dto.Layout, error) {
	return srv.layoutService.GetAvailableLayouts(ctx, userId)
}

func (srv *Service) DeleteLayout(ctx context.Context, req request.LayoutIdRequest, userId uuid.UUID) error {
	return srv.layoutService.DeleteLayoutById(ctx, req.LayoutId, userId)
}

func (srv *Service) UpdateLayout(ctx context.Context, req request.UpdateLayout, userId uuid.UUID) error {
	return srv.layoutService.UpdateLayout(ctx, req, userId)
}

func (srv *Service) ExportInfo(ctx context.Context, req dto.ExportInfoRequest) (*dto.ExportInfo, error) {
	return srv.layoutService.ExportLayouts(ctx, req.UserId)
}

func (srv *Service) ImportLayouts(ctx context.Context, userId uuid.UUID, req *dto.ImportInfoRequest) error {
	return srv.layoutService.ImportLayouts(ctx, userId, &req.Info)
}
