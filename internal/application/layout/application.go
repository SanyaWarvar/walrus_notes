package layout

import (
	"context"
	"wn/internal/domain/dto"
	"wn/internal/domain/events"
	"wn/pkg/applogger"
	"wn/pkg/trx"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type layoutService interface {
	CreateLayout(ctx context.Context, title, color string, ownerId uuid.UUID, isMain bool) (uuid.UUID, error)
	DeleteLayoutById(ctx context.Context, layoutId, ownerId uuid.UUID) error
	GetAvailableLayouts(ctx context.Context, userId uuid.UUID) ([]dto.Layout, error)
	ExportLayouts(ctx context.Context, userId uuid.UUID) (*dto.ExportInfo, error)
	UpdateLayout(ctx context.Context, req dto.UpdateLayout, userId uuid.UUID) error
	ImportLayouts(ctx context.Context, userId uuid.UUID, info *dto.ExportInfo) error
}

type permissionsService interface {
	CheckPermissionByLayoutId(ctx context.Context, targetId, userId uuid.UUID, read, write, edit bool) error
	GetAssociatedUsersByLayout(ctx context.Context, layoutId uuid.UUID) ([]uuid.UUID, error)
}

type eventProducer interface {
	SendToAssociatedUsers(ctx context.Context, targetId uuid.UUID, recipients []uuid.UUID, event events.Event) error
}

type Service struct {
	tx     trx.TransactionManager
	logger applogger.Logger

	layoutService      layoutService
	permissionsService permissionsService
	eventProducer      eventProducer
}

func NewService(
	tx trx.TransactionManager,
	logger applogger.Logger,
	layoutService layoutService,
	permissionsService permissionsService,
	eventProducer eventProducer,
) *Service {
	return &Service{
		tx:                 tx,
		logger:             logger,
		layoutService:      layoutService,
		permissionsService: permissionsService,
		eventProducer:      eventProducer,
	}
}

func (srv *Service) CreateLayout(ctx context.Context, req dto.NewLayoutRequest, userId uuid.UUID) (uuid.UUID, error) {
	id, err := srv.layoutService.CreateLayout(ctx, req.Title, req.Color, userId, false)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (srv *Service) GetLayoutsByUserId(ctx context.Context, userId uuid.UUID) ([]dto.Layout, error) {
	return srv.layoutService.GetAvailableLayouts(ctx, userId)
}

func (srv *Service) DeleteLayout(ctx context.Context, req dto.LayoutIdRequest, userId uuid.UUID) error {
	if err := srv.permissionsService.CheckPermissionByLayoutId(ctx, req.LayoutId, userId, true, false, true); err != nil {
		srv.logger.Warnf("DeleteLayout checkPerms: %s", err.Error())
		return err
	}

	recipients, err := srv.permissionsService.GetAssociatedUsersByLayout(ctx, req.LayoutId)
	if err != nil {
		return errors.Wrap(err, "p.permissionsService.GetAssociatedUsersByLayout")
	}

	err = srv.layoutService.DeleteLayoutById(ctx, req.LayoutId, userId)
	if err != nil {
		return err
	}

	go srv.eventProducer.SendToAssociatedUsers(context.Background(), req.LayoutId, recipients, &events.DeleteLayoutEvent{
		LayoutId: req.LayoutId,
	})

	return nil
}

func (srv *Service) UpdateLayout(ctx context.Context, req dto.UpdateLayout, userId uuid.UUID) error {
	if err := srv.permissionsService.CheckPermissionByLayoutId(ctx, req.LayoutId, userId, true, false, true); err != nil {
		srv.logger.Warnf("DeleteLayout checkPerms: %s", err.Error())
		return err
	}

	recipients, err := srv.permissionsService.GetAssociatedUsersByLayout(ctx, req.LayoutId)
	if err != nil {
		return errors.Wrap(err, "p.permissionsService.GetAssociatedUsersByLayout")
	}
	err = srv.layoutService.UpdateLayout(ctx, req, userId)
	if err != nil {
		return err
	}

	go srv.eventProducer.SendToAssociatedUsers(context.Background(), req.LayoutId, recipients, &events.UpdateLayoutEvent{
		LayoutId: req.LayoutId,
	})

	return nil
}

func (srv *Service) ExportInfo(ctx context.Context, req dto.ExportInfoRequest) (*dto.ExportInfo, error) {
	return srv.layoutService.ExportLayouts(ctx, req.UserId)
}

func (srv *Service) ImportLayouts(ctx context.Context, userId uuid.UUID, req *dto.ImportInfoRequest) error {
	return srv.layoutService.ImportLayouts(ctx, userId, &req.Info)
}
