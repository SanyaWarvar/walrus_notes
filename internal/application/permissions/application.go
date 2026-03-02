package permissions

import (
	"context"
	"fmt"
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
	SavePermissionsLink(ctx context.Context, item *dto.PermissionsToken, id uuid.UUID, ttl *time.Duration) error
	GetPermissionsLink(ctx context.Context, id uuid.UUID) (*dto.PermissionsToken, bool, error)
}

type noteRepository interface {
	GetByOwnerId(ctx context.Context, ownerId, noteId uuid.UUID) (*entity.Note, error)
}

type layoutRepository interface {
	GetByOwnerId(ctx context.Context, ownerId, layoutId uuid.UUID) (*entity.Layout, error)
}

type Application struct {
	tx     trx.TransactionManager
	logger applogger.Logger

	permissionsRepository     permissionsRepository
	permissionsLinkRepository permissionsLinkRepository
	layoutRepository          layoutRepository
	noteRepository            noteRepository
}

func NewApplication(
	tx trx.TransactionManager,
	logger applogger.Logger,
	permissionsRepository permissionsRepository,
	permissionsLinkRepository permissionsLinkRepository,
	layoutRepository layoutRepository,
	noteRepository noteRepository,
) *Application {
	return &Application{
		tx:                        tx,
		logger:                    logger,
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
	fmt.Println(ttl)

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

	if err = srv.permissionsLinkRepository.SavePermissionsLink(ctx, &dto.PermissionsToken{
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
