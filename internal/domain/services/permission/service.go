package permission

import (
	"wn/internal/domain/dto"
	"wn/internal/entity"
)

type Service struct{}

func NewPermissionsService() *Service {
	return &Service{}
}

func (srv *Service) ApplyUpdateRequest(req *dto.UpdatePermissionRequest, e *entity.Permission) *entity.Permission {
	e.CanRead = req.CanRead
	e.CanEdit = req.CanEdit
	e.CanWrite = req.CanWrite
	return e
}
