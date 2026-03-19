package permissions

import (
	"context"
	"time"
	"wn/internal/domain/dto"
	"wn/internal/domain/enum"
	"wn/internal/entity"
	apperrors "wn/internal/errors"
	"wn/pkg/apperror"
	"wn/pkg/applogger"
	"wn/pkg/constants"
	"wn/pkg/trx"
	"wn/pkg/util"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type permissionsRepository interface {
	GetPermission(ctx context.Context, filter *dto.GetPermissionsFilter) (*entity.Permission, error)
	GetPermissions(ctx context.Context, filter *dto.GetPermissionsFilter) ([]entity.Permission, error)
	DeletePermissions(ctx context.Context, permissionsIds ...uuid.UUID) error
	UpdatePermissions(ctx context.Context, item *entity.Permission) error
	CreatePermissions(ctx context.Context, item *entity.Permission) error
}

type permissionsLinkRepository interface {
	SavePermissionsLink(ctx context.Context, item *dto.PermissionToken, id uuid.UUID, ttl *time.Duration) error
	GetPermissionsLink(ctx context.Context, id uuid.UUID) (*dto.PermissionToken, bool, error)
}

type noteRepository interface {
	GetByOwnerId(ctx context.Context, ownerId, noteId uuid.UUID) (*entity.Note, error)
}

type layoutRepository interface {
	GetByOwnerId(ctx context.Context, ownerId, layoutId uuid.UUID) (*entity.Layout, error)
}

type permissionsService interface {
	ApplyUpdateRequest(req *dto.UpdatePermissionRequest, e *entity.Permission) *entity.Permission
}

type Application struct {
	tx     trx.TransactionManager
	logger applogger.Logger

	permissionsService permissionsService

	permissionsRepository     permissionsRepository
	permissionsLinkRepository permissionsLinkRepository
	layoutRepository          layoutRepository
	noteRepository            noteRepository
}

func NewApplication(
	tx trx.TransactionManager,
	logger applogger.Logger,
	permissionsService permissionsService,
	permissionsRepository permissionsRepository,
	permissionsLinkRepository permissionsLinkRepository,
	layoutRepository layoutRepository,
	noteRepository noteRepository,
) *Application {
	return &Application{
		tx:                        tx,
		logger:                    logger,
		permissionsService:        permissionsService,
		permissionsRepository:     permissionsRepository,
		permissionsLinkRepository: permissionsLinkRepository,
		layoutRepository:          layoutRepository,
		noteRepository:            noteRepository,
	}
}

func (srv *Application) GeneratePermissionsLink(ctx context.Context, userId uuid.UUID, req *dto.GeneratePermissionLinkRequest) (*dto.GeneratePermissionsLinkResponse, error) {
	id := uuid.New()
	ttl := req.ExpiredAt.Sub(util.GetCurrentUTCTime())
	if ttl < 0 {
		return nil, apperror.NewBadRequestError("bad expired at", constants.BindBodyError)
	}
	var err error

	switch req.Kind {
	case enum.PermissionsKindNote:
		_, err = srv.noteRepository.GetByOwnerId(ctx, userId, req.TargetId)
	case enum.PermissionsKindLayout:
		_, err = srv.layoutRepository.GetByOwnerId(ctx, userId, req.TargetId)
	default:
		return nil, apperrors.BadKind
	}

	if err != nil {
		return nil, apperrors.PermissionsNotEnough
	}

	if err = srv.permissionsLinkRepository.SavePermissionsLink(ctx, &dto.PermissionToken{
		FromUserId: userId,
		TargetId:   req.TargetId,
		Kind:       req.Kind,
		CanRead:    req.CanRead,
		CanWrite:   req.CanWrite,
		CanEdit:    req.CanEdit,
		ExpiredAt:  req.ExpiredAt,
	}, id, &ttl); err != nil {
		return nil, errors.Wrap(err, "SaveLink")
	}

	return &dto.GeneratePermissionsLinkResponse{
		LinkId: id,
	}, nil
}

func (srv *Application) ApplyPermissionsLink(ctx context.Context, userId uuid.UUID, req *dto.ApplyPermissionsRequest) error {

	perm, ex, err := srv.permissionsLinkRepository.GetPermissionsLink(ctx, req.LinkId)
	if err != nil {
		return errors.Wrap(err, "GetPermissionsLink")
	}

	if !ex {
		return apperrors.RecordNotFound
	}

	if userId == perm.FromUserId {
		return apperrors.CantApply
	}

	item, err := srv.permissionsRepository.GetPermission(ctx, &dto.GetPermissionsFilter{
		ToUserId: &userId,
		TargetId: &perm.TargetId,
	})
	if err != nil && err != apperrors.RecordNotFound {
		return errors.Wrap(err, "GetPermission")
	}

	if item != nil {
		return apperrors.AlreadyExist
	}

	return srv.permissionsRepository.CreatePermissions(ctx, &entity.Permission{
		Id:         uuid.New(),
		ToUserId:   userId,
		FromUserId: perm.FromUserId,
		TargetId:   perm.TargetId,
		Kind:       perm.Kind,
		CanRead:    perm.CanRead,
		CanWrite:   perm.CanWrite,
		CanEdit:    perm.CanEdit,
		CreatedAt:  util.GetCurrentUTCTime(),
	})
}

func (srv *Application) GetPermissionsDashboard(ctx context.Context, userId uuid.UUID) (*dto.PermissionsDashbord, error) {
	recivied, err := srv.permissionsRepository.GetPermissions(ctx, &dto.GetPermissionsFilter{
		ToUserId: &userId,
	})
	if err != nil {
		return nil, errors.Wrap(err, "recivied")
	}

	shared, err := srv.permissionsRepository.GetPermissions(ctx, &dto.GetPermissionsFilter{
		FromUserId: &userId,
	})
	if err != nil {
		return nil, errors.Wrap(err, "shared")
	}

	sharedDto := make([]dto.Permission, 0, len(shared))
	for i := range shared {
		sharedDto = append(sharedDto, *dto.PermissionFromEntity(&shared[i]))
	}

	reciviedDto := make([]dto.Permission, 0, len(recivied))
	for i := range recivied {
		reciviedDto = append(sharedDto, *dto.PermissionFromEntity(&recivied[i]))
	}

	return &dto.PermissionsDashbord{
		Shared:   sharedDto,
		Recivied: reciviedDto,
	}, nil
}

func (srv *Application) DeletePermission(ctx context.Context, userId uuid.UUID, req *dto.DeletePermissionsRequest) error {
	permission, err := srv.permissionsRepository.GetPermission(ctx, &dto.GetPermissionsFilter{
		Id: &req.PermissionId,
	})
	if err != nil {
		return err
	}
	if permission.FromUserId != userId || permission.ToUserId != userId {
		return apperrors.PermissionsNotEnough
	}
	return srv.permissionsRepository.DeletePermissions(ctx, req.PermissionId)
}

func (srv *Application) UpdatePermission(ctx context.Context, userId uuid.UUID, req *dto.UpdatePermissionRequest) error {
	permission, err := srv.permissionsRepository.GetPermission(ctx, &dto.GetPermissionsFilter{
		Id: &req.PermissionId,
	})
	if err != nil {
		return err
	}
	if permission.FromUserId != userId {
		return apperrors.PermissionsNotEnough
	}

	permission = srv.permissionsService.ApplyUpdateRequest(req, permission)

	return srv.permissionsRepository.UpdatePermissions(ctx, permission)
}
