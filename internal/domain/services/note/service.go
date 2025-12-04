package note

import (
	"context"
	"encoding/json"
	"math"
	"sort"
	"wn/internal/domain/dto"
	"wn/internal/entity"
	"wn/pkg/applogger"
	"wn/pkg/trx"
	"wn/pkg/util"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type noteRepo interface {
	DeleteNoteById(ctx context.Context, noteId, userId uuid.UUID) error
	CreateNote(ctx context.Context, item *entity.Note) (uuid.UUID, error)
	UpdateNote(ctx context.Context, newItem *entity.Note) error
	GetNoteCountInLayout(ctx context.Context, layoutId uuid.UUID) (int, error)
	GetNotesByLayoutId(ctx context.Context, layoutId, userId uuid.UUID, offset, limit int) ([]entity.Note, error)
	GetNotesWithPosition(ctx context.Context, layoutId, userId uuid.UUID) ([]entity.NoteWithPosition, error)
	GetNotesWithoutPosition(ctx context.Context, layoutId, userId uuid.UUID) ([]entity.Note, error)
	SearchNotes(ctx context.Context, userId uuid.UUID, search string) ([]entity.Note, error)
	UpdateDraftById(ctx context.Context, userId, noteId uuid.UUID, newDraft string) error
	CommitDraft(ctx context.Context, userId, noteId uuid.UUID) error
}

type positionsRepo interface {
	CreateNotePosition(ctx context.Context, noteId uuid.UUID, xPos, yPos *float64) error
	UpdateNotePosition(ctx context.Context, noteId uuid.UUID, xPos, yPos *float64) error
	DeleteNotesPositionByNoteId(ctx context.Context, noteId uuid.UUID) error
}

type linksRepo interface {
	DeleteLinksWithNote(ctx context.Context, noteId uuid.UUID) error
	DeleteLink(ctx context.Context, firstNoteId, secondNoteId uuid.UUID) error
	LinkNotes(ctx context.Context, firstNoteId, secondNoteId uuid.UUID) error
	GetAllLinks(ctx context.Context, noteIds []uuid.UUID) ([]entity.Link, error)
}

type layoutRepo interface {
}

type Service struct {
	tx     trx.TransactionManager
	logger applogger.Logger

	noteRepo      noteRepo
	layoutRepo    layoutRepo
	linksRepo     linksRepo
	positionsRepo positionsRepo
}

func NewService(
	tx trx.TransactionManager,
	logger applogger.Logger,
	noteRepo noteRepo,
	layoutRepo layoutRepo,
	linksRepo linksRepo,
	positionsRepo positionsRepo,
) *Service {
	return &Service{
		tx:            tx,
		logger:        logger,
		noteRepo:      noteRepo,
		layoutRepo:    layoutRepo,
		linksRepo:     linksRepo,
		positionsRepo: positionsRepo,
	}
}

func (srv *Service) DeleteLink(ctx context.Context, noteId1, noteId2 uuid.UUID) error {
	return srv.linksRepo.DeleteLink(ctx, noteId1, noteId2)
}

func (srv *Service) DeleteNoteById(ctx context.Context, noteId, userId uuid.UUID) error {
	return srv.tx.Transaction(ctx, func(ctx context.Context) error {
		err := srv.positionsRepo.DeleteNotesPositionByNoteId(ctx, noteId)
		if err != nil {
			return err
		}
		err = srv.linksRepo.DeleteLinksWithNote(ctx, noteId)
		if err != nil {
			return err
		}
		err = srv.noteRepo.DeleteNoteById(ctx, noteId, userId)
		if err != nil {
			return err
		}
		return nil
	})
}

func (srv *Service) CreateNote(ctx context.Context, title, payload string, ownerId, layoutId, mainLayoutId uuid.UUID) (uuid.UUID, error) {
	n := entity.Note{
		Id:         util.NewUUID(),
		Title:      title,
		Payload:    payload,
		CreatedAt:  util.GetCurrentUTCTime(),
		OwnerId:    ownerId,
		HaveAccess: []uuid.UUID{ownerId},
		LayoutId:   layoutId,
	}

	return n.Id, srv.tx.Transaction(ctx, func(ctx context.Context) error {
		_, err := srv.noteRepo.CreateNote(ctx, &n)
		if err != nil {
			return errors.Wrap(err, "srv.noteRepo.CreateNote")
		}
		err = srv.positionsRepo.CreateNotePosition(ctx, n.Id, nil, nil)
		if err != nil {
			return errors.Wrap(err, "srv.noteRepo.CreateNote")
		}
		return nil
	})
}

func (srv *Service) UpdateNote(ctx context.Context, title, payload string, noteId, ownerId uuid.UUID) error {
	n := entity.Note{
		Id:      noteId,
		Title:   title,
		Payload: payload,
		OwnerId: ownerId,
	}
	return srv.noteRepo.UpdateNote(ctx, &n)
}

func (srv *Service) GetNotesWithPagination(ctx context.Context, page int, layoutId, userId uuid.UUID) ([]dto.Note, int, error) {

	count, err := srv.noteRepo.GetNoteCountInLayout(ctx, layoutId)
	if err != nil {
		return nil, 0, errors.Wrap(err, "srv.noteRepo.GetNoteCountInLayout")
	}
	offset := util.CalculateOffset(page)
	limit := util.CalculateLimit()
	notes, err := srv.noteRepo.GetNotesByLayoutId(ctx, layoutId, userId, offset, limit)
	if err != nil {
		return nil, 0, errors.Wrap(err, "srv.noteRepo.GetNotesByLayoutId")
	}
	links, err := srv.linksRepo.GetAllLinks(ctx, getIds(notes))
	notesDto := dto.NotesFromEntities(notes, links)
	return notesDto, count, err
}

func (srv *Service) GetNotesWithoutPosition(ctx context.Context, layoutId, userId uuid.UUID) ([]dto.Note, error) {
	notes, err := srv.noteRepo.GetNotesWithoutPosition(ctx, layoutId, userId)
	if err != nil {
		return nil, err
	}
	links, err := srv.linksRepo.GetAllLinks(ctx, getIds(notes))
	notesDto := dto.NotesFromEntities(notes, links)
	return notesDto, err
}

func (srv *Service) GetNotesWithPosition(ctx context.Context, mainLayoutId, layoutId, userId uuid.UUID) ([]dto.Note, error) {
	t := layoutId
	if mainLayoutId == t {
		t = uuid.Nil
	}
	notes, err := srv.noteRepo.GetNotesWithPosition(ctx, t, userId)
	if err != nil {
		return nil, err
	}
	links, err := srv.linksRepo.GetAllLinks(ctx, getIds(notes))
	notesDto := dto.NotesFromEntitiesWithPosition(notes, links)
	return notesDto, err
}

func (srv *Service) UpdateNotePosition(ctx context.Context, noteId uuid.UUID, xPos, yPos *float64) error {
	return srv.positionsRepo.UpdateNotePosition(ctx, noteId, xPos, yPos)
}

func (srv *Service) CreateLink(ctx context.Context, noteId1, noteId2 uuid.UUID) error {
	return srv.linksRepo.LinkNotes(ctx, noteId1, noteId2)
}

// todo добавлять беклинки?
func (srv *Service) SearchNotes(ctx context.Context, userId uuid.UUID, search string) ([]dto.Note, error) {
	notes, err := srv.noteRepo.SearchNotes(ctx, userId, search)
	if err != nil {
		return nil, err
	}

	notesDto := dto.NotesFromEntities(notes, nil)
	return notesDto, err
}

func (srv *Service) HandleCreateDraft(msg *dto.SocketMessage, userId uuid.UUID) (*dto.SocketMessage, error) {
	ctx := context.Background()
	var item dto.DraftNote
	err := json.Unmarshal(msg.Payload, &item)
	if err != nil {
		return &dto.SocketMessage{
			Event:   "COMMIT_DRAFT_RESPONSE",
			Payload: []byte("{\"status\": \"false\"}"),
		}, err
	}
	err = srv.noteRepo.UpdateDraftById(ctx, userId, item.NoteId, item.NewDraft)
	return &dto.SocketMessage{
		Event:   "UPDATE_DRAFT_RESPONSE",
		Payload: []byte("{\"status\": \"true\"}"),
	}, nil
}

func (srv *Service) HandleCommitDraft(msg *dto.SocketMessage, userId uuid.UUID) (*dto.SocketMessage, error) {
	ctx := context.Background()
	var item dto.CommitDraftNote
	err := json.Unmarshal(msg.Payload, &item)
	if err != nil {
		return &dto.SocketMessage{
			Event:   "COMMIT_DRAFT_RESPONSE",
			Payload: []byte("{\"status\": \"false\"}"),
		}, err
	}
	err = srv.noteRepo.CommitDraft(ctx, userId, item.NoteId)
	return &dto.SocketMessage{
		Event:   "COMMIT_DRAFT_RESPONSE",
		Payload: []byte("{\"status\": \"true\"}"),
	}, nil
}

func (srv *Service) GenerateCluster(notes []dto.Note) []dto.Note {
	if len(notes) == 0 {
		return notes
	}

	// Группируем заметки по layoutId
	layoutGroups := make(map[uuid.UUID][]dto.Note)
	for _, note := range notes {
		layoutGroups[note.LayoutId] = append(layoutGroups[note.LayoutId], note)
	}

	// Параметры кластеризации
	clusterSpacing := 100.0 // расстояние между кластерами
	gridSpacing := 20.0     // расстояние между заметками внутри кластера
	notesPerRow := 5        // максимальное количество заметок в строке кластера

	// Обрабатываем каждый кластер
	clusterX := 0.0
	result := make([]dto.Note, 0, len(notes))

	for _, groupNotes := range layoutGroups {
		// Сортируем заметки внутри кластера для детерминированного позиционирования
		sort.Slice(groupNotes, func(i, j int) bool {
			return groupNotes[i].Id.String() < groupNotes[j].Id.String()
		})

		// Позиционируем заметки внутри кластера
		for i, note := range groupNotes {
			// Если у заметки уже есть позиция, используем её (относительно кластера)
			if note.Position != nil {
				note.Position.XPos += clusterX
				result = append(result, note)
				continue
			}

			// Автоматическое позиционирование внутри кластера
			row := i / notesPerRow
			col := i % notesPerRow

			xPos := clusterX + float64(col)*gridSpacing
			yPos := float64(row) * gridSpacing

			// Создаем новую позицию
			note.Position = &dto.Position{
				XPos: xPos,
				YPos: yPos,
			}
			result = append(result, note)
		}

		// Сдвигаем следующий кластер вправо
		maxNotesInCluster := len(groupNotes)
		clusterWidth := math.Min(float64(notesPerRow), float64(maxNotesInCluster)) * gridSpacing
		clusterX += clusterWidth + clusterSpacing
	}

	return result
}

func calculateClusterBounds(notes []dto.Note) (minX, minY, maxX, maxY float64) {
	if len(notes) == 0 {
		return 0, 0, 0, 0
	}

	// Инициализируем первыми валидными координатами
	for _, note := range notes {
		if note.Position != nil {
			minX = note.Position.XPos
			minY = note.Position.YPos
			maxX = note.Position.XPos
			maxY = note.Position.YPos
			break
		}
	}

	// Находим границы кластера
	for _, note := range notes {
		if note.Position != nil {
			if note.Position.XPos < minX {
				minX = note.Position.XPos
			}
			if note.Position.YPos < minY {
				minY = note.Position.YPos
			}
			if note.Position.XPos > maxX {
				maxX = note.Position.XPos
			}
			if note.Position.YPos > maxY {
				maxY = note.Position.YPos
			}
		}
	}

	// Защита от вырожденного случая (все точки в одном месте)
	if maxX == minX {
		maxX = minX + 1
	}
	if maxY == minY {
		maxY = minY + 1
	}

	return minX, minY, maxX, maxY
}
