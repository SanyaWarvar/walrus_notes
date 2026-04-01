package note

import (
	"context"
	"wn/internal/domain/dto"
	req "wn/internal/domain/dto/request"
	"wn/pkg/applogger"
	"wn/pkg/trx"

	"github.com/google/uuid"
)

type noteService interface {
	DeleteNoteById(ctx context.Context, noteId uuid.UUID) error
	CreateNote(ctx context.Context, title, payload string, ownerId, layoutId, mainLayoutId uuid.UUID) (uuid.UUID, error)
	UpdateNote(ctx context.Context, title, payload string, noteId uuid.UUID) error
	GetNotesWithPagination(ctx context.Context, page int, layoutId, userId uuid.UUID) ([]dto.Note, int, error)
	GetNotesWithPosition(ctx context.Context, mainLayoutId, layoutId, userId uuid.UUID) ([]dto.Note, error)
	GetNotesWithoutPosition(ctx context.Context, layoutId, userId uuid.UUID) ([]dto.Note, error)
	UpdateNotePosition(ctx context.Context, noteId uuid.UUID, xPos, yPos *float64) error
	CreateLink(ctx context.Context, noteId1, noteId2 uuid.UUID) error
	DeleteLink(ctx context.Context, noteId1, noteId2 uuid.UUID) error
	SearchNotes(ctx context.Context, userId uuid.UUID, search string) ([]dto.Note, error)
	GenerateCluster(notes []dto.Note) []dto.Note
	DragNote(ctx context.Context, noteId, toLayout uuid.UUID) error
}

type permissionsService interface {
	CheckPermissionByLayoutId(ctx context.Context, targetId, userId uuid.UUID, read, write, edit bool) error
	CheckPermissionByNoteId(ctx context.Context, targetId, userId uuid.UUID, read, write, edit bool) error
}

type Service struct {
	tx     trx.TransactionManager
	logger applogger.Logger

	noteService        noteService
	permissionsService permissionsService
}

func NewService(
	tx trx.TransactionManager,
	logger applogger.Logger,
	noteService noteService,
	permissionsService permissionsService,
) *Service {
	return &Service{
		tx:                 tx,
		logger:             logger,
		noteService:        noteService,
		permissionsService: permissionsService,
	}
}

func (srv *Service) CreateNote(ctx context.Context, req req.NoteRequest, userId uuid.UUID, mainLayoutId uuid.UUID) (uuid.UUID, error) {
	if err := srv.permissionsService.CheckPermissionByLayoutId(ctx, req.LayoutId, userId, true, true, false); err != nil {
		srv.logger.Warnf("CreateNote checkPerms: %s", err.Error())
		return uuid.Nil, err
	}
	return srv.noteService.CreateNote(ctx, req.Title, req.Payload, userId, req.LayoutId, mainLayoutId)
}

func (srv *Service) UpdateNote(ctx context.Context, req req.NoteWithIdRequest, userId uuid.UUID) error {
	if err := srv.permissionsService.CheckPermissionByNoteId(ctx, req.NoteId, userId, true, true, false); err != nil {
		srv.logger.Warnf("UpdateNote checkPerms: %s", err.Error())
		return err
	}

	return srv.noteService.UpdateNote(ctx, req.Title, req.Payload, req.NoteId)
}

func (srv *Service) DeleteNote(ctx context.Context, req req.NoteId, userId, mainLayoutId uuid.UUID) error {
	if err := srv.permissionsService.CheckPermissionByNoteId(ctx, req.NoteId, userId, true, true, false); err != nil {
		srv.logger.Warnf("DeleteNote checkPerms: %s", err.Error())
		return err
	}

	return srv.noteService.DeleteNoteById(ctx, req.NoteId)
}

func (srv *Service) GetNotesFromLayout(ctx context.Context, req req.GetNotesFromLayoutRequest, userId uuid.UUID) ([]dto.Note, int, error) {
	if err := srv.permissionsService.CheckPermissionByLayoutId(ctx, req.LayoutId, userId, true, false, false); err != nil {
		srv.logger.Warnf("GetNotesFromLayout checkPerms: %s", err.Error())
		return nil, 0, err
	}
	return srv.noteService.GetNotesWithPagination(ctx, req.Page, req.LayoutId, userId)
}

func (srv *Service) GetNotesWithPosition(ctx context.Context, userId, mainLayoutId uuid.UUID, req req.GetNotesFromLayoutWithoutPagRequest) ([]dto.Note, error) {
	if err := srv.permissionsService.CheckPermissionByLayoutId(ctx, req.LayoutId, userId, true, false, false); err != nil {
		srv.logger.Warnf("GetNotesFromLayout checkPerms: %s", err.Error())
		return nil, err
	}
	notes, err := srv.noteService.GetNotesWithPosition(ctx, mainLayoutId, req.LayoutId, userId)
	if err != nil {
		return nil, err
	}
	if mainLayoutId == req.LayoutId {
		return srv.noteService.GenerateCluster(notes), nil
	}
	return notes, err
}

func (srv *Service) GetNotesWithoutPosition(ctx context.Context, userId uuid.UUID, req req.GetNotesFromLayoutWithoutPagRequest) ([]dto.Note, error) {
	if err := srv.permissionsService.CheckPermissionByLayoutId(ctx, req.LayoutId, userId, true, false, false); err != nil {
		srv.logger.Warnf("GetNotesFromLayout checkPerms: %s", err.Error())
		return nil, err
	}
	return srv.noteService.GetNotesWithoutPosition(ctx, req.LayoutId, userId)
}

func (srv *Service) UpdateNotePosition(ctx context.Context, userId uuid.UUID, req req.UpdateNotePositionRequest) error {
	if err := srv.permissionsService.CheckPermissionByLayoutId(ctx, req.LayoutId, userId, true, true, false); err != nil {
		srv.logger.Warnf("GetNotesFromLayout checkPerms: %s", err.Error())
		return err
	}
	return srv.noteService.UpdateNotePosition(ctx, req.NoteId, req.XPos, req.YPos)
}

func (srv *Service) CreateLink(ctx context.Context, userId uuid.UUID, req req.LinkBetweenNotesRequest) error {
	if err := srv.permissionsService.CheckPermissionByLayoutId(ctx, req.LayoutId, userId, true, true, false); err != nil {
		srv.logger.Warnf("GetNotesFromLayout checkPerms: %s", err.Error())
		return err
	}
	return srv.noteService.CreateLink(ctx, req.FirstNoteId, req.SecondNoteId)
}

func (srv *Service) DeleteLink(ctx context.Context, userId uuid.UUID, req req.LinkBetweenNotesRequest) error {
	if err := srv.permissionsService.CheckPermissionByLayoutId(ctx, req.LayoutId, userId, true, true, false); err != nil {
		srv.logger.Warnf("GetNotesFromLayout checkPerms: %s", err.Error())
		return err
	}
	return srv.noteService.DeleteLink(ctx, req.FirstNoteId, req.SecondNoteId)
}

func (srv *Service) DragNote(ctx context.Context, userId uuid.UUID, req req.DragNoteRequest) error {
	if err := srv.permissionsService.CheckPermissionByNoteId(ctx, req.NoteId, userId, true, true, false); err != nil {
		srv.logger.Warnf("CheckPermissionByNoteId: %s", err.Error())
		return err
	}
	if err := srv.permissionsService.CheckPermissionByLayoutId(ctx, req.ToLayoutId, userId, true, true, false); err != nil {
		srv.logger.Warnf("CheckPermissionByLayoutId: %s", err.Error())
		return err
	}

	return srv.noteService.DragNote(ctx, req.NoteId, req.ToLayoutId)
}

func (srv *Service) SearchNotes(ctx context.Context, userId uuid.UUID, search string) ([]dto.Note, error) {
	return srv.noteService.SearchNotes(ctx, userId, search)
}
