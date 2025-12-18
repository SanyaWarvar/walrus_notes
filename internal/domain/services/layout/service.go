package layout

import (
	"context"
	"fmt"
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
	DeleteLinksByLayoutId(ctx context.Context, layoutId uuid.UUID) error
	LinkNotes(ctx context.Context, firstNoteId, secondNoteId uuid.UUID) error
	GetAllLinks(ctx context.Context, noteIds []uuid.UUID) ([]entity.Link, error)
}

type noteRepo interface {
	DeleteNotesByLayoutId(ctx context.Context, layoutId, userId uuid.UUID) error
	GetNotesWithPosition(ctx context.Context, layoutId, userId uuid.UUID) ([]entity.NoteWithPosition, error)
	GetFullNotesByLayoutId(ctx context.Context, layoutId, userId uuid.UUID) ([]dto.Note, error)
}

type noteService interface {
	RessurectNotes(ctx context.Context, item *dto.Note) error
}

type positionsRepo interface {
	DeleteNotesPositionsByLayoutId(ctx context.Context, layoutId uuid.UUID) error
}

type Service struct {
	tx     trx.TransactionManager
	logger applogger.Logger

	layoutRepo    layoutRepo
	linksRepo     linksRepo
	noteRepo      noteRepo
	positionsRepo positionsRepo
	noteService   noteService
}

func NewService(
	tx trx.TransactionManager,
	logger applogger.Logger,
	layoutRepo layoutRepo,
	linksRepo linksRepo,
	noteRepo noteRepo,
	positionsRepo positionsRepo,
	noteService noteService,
) *Service {
	return &Service{
		tx:            tx,
		logger:        logger,
		layoutRepo:    layoutRepo,
		linksRepo:     linksRepo,
		noteRepo:      noteRepo,
		noteService:   noteService,
		positionsRepo: positionsRepo,
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

		err := srv.positionsRepo.DeleteNotesPositionsByLayoutId(ctx, layoutId)
		if err != nil {
			return errors.Wrap(err, "srv.positionsRepo.DeleteNotesPositionsByLayoutId")
		}

		err = srv.linksRepo.DeleteLinksByLayoutId(ctx, layoutId)
		if err != nil {
			return errors.Wrap(err, "srv.linksRepo.DeleteLinksByLayoutId")
		}

		err = srv.noteRepo.DeleteNotesByLayoutId(ctx, layoutId, ownerId)
		if err != nil {
			return errors.Wrap(err, "srv.noteRepo.DeleteNotesByLayoutId")
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

func (srv *Service) ExportLayouts(ctx context.Context, userId uuid.UUID) (*dto.ExportInfo, error) {
	layouts, err := srv.layoutRepo.GetAvailableLayouts(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "GetAvailableLayouts")
	}
	var output dto.ExportInfo
	output.UserId = userId
	output.CreatedAt = util.GetCurrentUTCTime()
	output.Notes = map[uuid.UUID][]dto.Note{}

	for _, l := range layouts {
		notes, err := srv.noteRepo.GetFullNotesByLayoutId(ctx, l.Id, userId)
		if err != nil {
			return nil, errors.Wrap(err, "GetFullNotesByLayoutId")
		}
		output.Layouts = append(output.Layouts, dto.Layout{
			Id:      l.Id,
			Title:   l.Title,
			OwnerId: l.OwnerId,
			IsMain:  l.IsMain,
			Color:   l.Color,
		})
		output.Notes[l.Id] = notes
	}
	return &output, nil
}

func (srv *Service) ImportLayouts(ctx context.Context, userId uuid.UUID, info *dto.ExportInfo) error {
	layouts, err := srv.layoutRepo.GetAvailableLayouts(ctx, userId)
	return srv.tx.Transaction(ctx, func(ctx context.Context) error {
		for i := range layouts {
			if layouts[i].IsMain {
				continue
			}
			err = srv.DeleteLayoutById(ctx, layouts[i].Id, userId)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("ImportLayouts: cant delete: %s for user %s", layouts[i].Id.String(), userId.String()))

			}
		}

		for _, l := range info.Layouts {
			if l.IsMain {
				continue
			}
			_, err := srv.layoutRepo.CreateLayout(ctx, &entity.Layout{
				Id:         l.Id,
				Title:      l.Title,
				Color:      l.Color,
				OwnerId:    l.OwnerId,
				HaveAccess: []uuid.UUID{l.OwnerId},
				IsMain:     false,
			})
			if err != nil {
				return errors.Wrap(err, "CreateLayout")
			}
			for _, note := range info.Notes[l.Id] {
				err = srv.noteService.RessurectNotes(ctx, &note)
				if err != nil {
					return errors.Wrap(err, "RessurectNotes")
				}
			}
		}
		for _, l := range info.Layouts {
			if l.IsMain {
				continue
			}
			for _, item := range info.Notes[l.Id] {
				for _, out := range item.LinkedWithOut {
					err = srv.linksRepo.LinkNotes(ctx, item.Id, out)
					if err != nil {
						return errors.Wrap(err, "LinkNotes")
					}
				}
			}
		}
		return nil
	})
}
