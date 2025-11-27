package note

import (
	"context"
	"encoding/json"
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
	UpdateNotePosition(ctx context.Context, layoutId, noteId uuid.UUID, xPos, yPos *float64) error
	SearchNotes(ctx context.Context, userId uuid.UUID, search string) ([]entity.Note, error)
	UpdateDraftById(ctx context.Context, userId, noteId uuid.UUID, newDraft string) error
	CommitDraft(ctx context.Context, userId, noteId uuid.UUID) error
}

type linksRepo interface {
	DeleteLinkNotes(ctx context.Context, layoutId, firstNoteId, secondNoteId uuid.UUID) error
	DeleteLinksFromLayout(ctx context.Context, layoutId uuid.UUID) error
	DeleteLinksWithNote(ctx context.Context, noteId uuid.UUID) error
	LinkNotes(ctx context.Context, layoutId, firstNoteId, secondNoteId uuid.UUID) error
	GetAllLinks(ctx context.Context, layoutId uuid.UUID, noteIds []uuid.UUID) ([]entity.Link, error)
}

type layoutRepo interface {
}

type Service struct {
	tx     trx.TransactionManager
	logger applogger.Logger

	noteRepo   noteRepo
	layoutRepo layoutRepo
	linksRepo  linksRepo
}

func NewService(
	tx trx.TransactionManager,
	logger applogger.Logger,
	noteRepo noteRepo,
	layoutRepo layoutRepo,
	linksRepo linksRepo,
) *Service {
	return &Service{
		tx:         tx,
		logger:     logger,
		noteRepo:   noteRepo,
		layoutRepo: layoutRepo,
		linksRepo:  linksRepo,
	}
}

func (srv *Service) DeleteNoteById(ctx context.Context, noteId, userId, mainLayoutId uuid.UUID) error {
	return srv.tx.Transaction(ctx, func(ctx context.Context) error {

		err := srv.noteRepo.DeleteNoteById(ctx, noteId, userId)
		if err != nil {
			return err
		}

		err = srv.linksRepo.DeleteLinksWithNote(ctx, noteId)
		if err != nil {
			return err
		}

		return nil
	})
}

// todo add check access
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
	id, err := srv.noteRepo.CreateNote(ctx, &n)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "srv.noteRepo.CreateNote")
	}
	return id, err
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
	links, err := srv.linksRepo.GetAllLinks(ctx, layoutId, getIds(notes))
	notesDto := dto.NotesFromEntities(notes, links)
	return notesDto, count, err
}

func (srv *Service) GetNotesWithoutPosition(ctx context.Context, layoutId, userId uuid.UUID) ([]dto.Note, error) {
	notes, err := srv.noteRepo.GetNotesWithoutPosition(ctx, layoutId, userId)
	if err != nil {
		return nil, err
	}
	links, err := srv.linksRepo.GetAllLinks(ctx, layoutId, getIds(notes))
	notesDto := dto.NotesFromEntities(notes, links)
	return notesDto, err
}

func (srv *Service) GetNotesWithPosition(ctx context.Context, layoutId, userId uuid.UUID) ([]dto.Note, error) {
	notes, err := srv.noteRepo.GetNotesWithPosition(ctx, layoutId, userId)
	if err != nil {
		return nil, err
	}
	links, err := srv.linksRepo.GetAllLinks(ctx, layoutId, getIds(notes))
	notesDto := dto.NotesFromEntitiesWithPosition(notes, links)
	return notesDto, err
}

func (srv *Service) UpdateNotePosition(ctx context.Context, layoutId, noteId uuid.UUID, xPos, yPos *float64) error {
	return srv.noteRepo.UpdateNotePosition(ctx, layoutId, noteId, xPos, yPos)
}

func (srv *Service) CreateLink(ctx context.Context, layoutId, noteId1, noteId2 uuid.UUID) error {
	return srv.linksRepo.LinkNotes(ctx, layoutId, noteId1, noteId2)
}

func (srv *Service) DeleteLink(ctx context.Context, layoutId, noteId1, noteId2 uuid.UUID) error {
	return srv.linksRepo.DeleteLinkNotes(ctx, layoutId, noteId1, noteId2)
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
	// Группируем заметки по layoutId
	clusters := map[uuid.UUID][]dto.Note{}
	for i := range notes {
		clusters[notes[i].LayoutId] = append(clusters[notes[i].LayoutId], notes[i])
	}

	var result []dto.Note

	for _, clusterNotes := range clusters {
		if len(clusterNotes) == 0 {
			continue
		}

		// Вычисляем bounding box для кластера
		minX, minY, maxX, maxY := calculateClusterBounds(clusterNotes)

		// Создаем смещение для этого кластера
		clusterOffsetX := minX
		clusterOffsetY := minY
		clusterWidth := maxX - minX
		clusterHeight := maxY - minY

		// Нормализуем координаты относительно кластера
		for i := range clusterNotes {
			if clusterNotes[i].Position != nil {
				// Преобразуем глобальные координаты в локальные (0-1 относительно кластера)
				localX := (clusterNotes[i].Position.XPos - clusterOffsetX) / clusterWidth
				localY := (clusterNotes[i].Position.YPos - clusterOffsetY) / clusterHeight

				// Обновляем позицию (можно также масштабировать если нужно)
				clusterNotes[i].Position.XPos = localX
				clusterNotes[i].Position.YPos = localY
			}

			result = append(result, clusterNotes[i])
		}
	}

	return result
}

// Альтернативный вариант - расположить кластеры в сетке
func (srv *Service) GenerateClusterGrid(notes []dto.Note, gridCols int) []dto.Note {
	clusters := map[uuid.UUID][]dto.Note{}
	for i := range notes {
		clusters[notes[i].LayoutId] = append(clusters[notes[i].LayoutId], notes[i])
	}

	var result []dto.Note
	clusterIndex := 0

	for _, clusterNotes := range clusters {
		if len(clusterNotes) == 0 {
			continue
		}

		// Вычисляем позицию кластера в сетке
		row := clusterIndex / gridCols
		col := clusterIndex % gridCols
		clusterBaseX := float64(col) * 1000 // смещение по X для кластера
		clusterBaseY := float64(row) * 1000 // смещение по Y для кластера

		// Вычисляем bounding box для нормализации
		minX, minY, maxX, maxY := calculateClusterBounds(clusterNotes)
		clusterWidth := maxX - minX
		clusterHeight := maxY - minY

		// Нормализуем и смещаем координаты
		for i := range clusterNotes {
			if clusterNotes[i].Position != nil {
				// Нормализуем к диапазону 0-1
				normalizedX := (clusterNotes[i].Position.XPos - minX) / clusterWidth
				normalizedY := (clusterNotes[i].Position.YPos - minY) / clusterHeight

				// Масштабируем и смещаем в позицию кластера
				// (например, каждая ячейка сетки 800x600)
				clusterNotes[i].Position.XPos = clusterBaseX + normalizedX*800
				clusterNotes[i].Position.YPos = clusterBaseY + normalizedY*600
			}

			result = append(result, clusterNotes[i])
		}

		clusterIndex++
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
