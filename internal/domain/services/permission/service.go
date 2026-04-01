package permission

import (
	"context"
	"wn/internal/domain/dto"
	"wn/internal/entity"
	apperrors "wn/internal/errors"

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

type layoutRepo interface {
	GetById(ctx context.Context, layoutId uuid.UUID) (*entity.Layout, error)
}

type noteRepo interface {
	GetById(ctx context.Context, noteId uuid.UUID) (*entity.Note, error)
}

type Service struct {
	permissionsRepository permissionsRepository
	layoutRepo            layoutRepo
	noteRepo              noteRepo
}

func NewPermissionsService(permissionsRepository permissionsRepository, layoutRepo layoutRepo, noteRepo noteRepo) *Service {
	return &Service{
		permissionsRepository: permissionsRepository,
		noteRepo:              noteRepo,
		layoutRepo:            layoutRepo,
	}
}

func (srv *Service) ApplyUpdateRequest(req *dto.UpdatePermissionRequest, e *entity.Permission) *entity.Permission {
	e.CanRead = req.CanRead
	e.CanEdit = req.CanEdit
	e.CanWrite = req.CanWrite
	return e
}

func (srv *Service) CheckPermissionByLayoutId(ctx context.Context, targetId, userId uuid.UUID, read, write, edit bool) error {
	l, err := srv.layoutRepo.GetById(ctx, targetId)
	if err != nil {
		return errors.Wrap(err, "srv.layoutRepo.GetById")
	}
	if l != nil && l.OwnerId == userId {
		return nil
	}

	perm, err := srv.permissionsRepository.GetPermission(ctx, &dto.GetPermissionsFilter{
		ToUserId: &userId,
		TargetId: &targetId,
	})
	if err != nil {
		return errors.Wrap(err, "GetPermission")
	}

	if edit && !perm.CanEdit {
		return apperrors.PermissionsNotEnough
	}
	if write && !perm.CanWrite {
		return apperrors.PermissionsNotEnough
	}
	if read && !perm.CanRead {
		return apperrors.PermissionsNotEnough
	}

	return nil
}

func (srv *Service) CheckPermissionByNoteId(ctx context.Context, targetId, userId uuid.UUID, read, write, edit bool) error {
	n, err := srv.noteRepo.GetById(ctx, targetId)
	if err != nil {
		return errors.Wrap(err, "noteRepo.GetById")
	}
	if n.OwnerId == userId {
		return nil
	}

	l, err := srv.layoutRepo.GetById(ctx, n.LayoutId)
	if err != nil {
		return errors.Wrap(err, "layoutRepo.GetById")
	}
	if l.OwnerId == userId {
		return nil
	}

	perm, err := srv.permissionsRepository.GetPermission(ctx, &dto.GetPermissionsFilter{
		ToUserId: &userId,
		TargetId: &n.LayoutId,
	})
	if err != nil {
		return errors.Wrap(err, "GetPermission")
	}

	if edit && !perm.CanEdit {
		return apperrors.PermissionsNotEnough
	}
	if write && !perm.CanWrite {
		return apperrors.PermissionsNotEnough
	}
	if read && !perm.CanRead {
		return apperrors.PermissionsNotEnough
	}

	return nil
}
