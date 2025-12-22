package note

import (
	"context"
	"wn/internal/domain/dto"
	"wn/internal/entity"
	apperrors "wn/internal/errors"
	"wn/internal/infrastructure/repository/common"
	"wn/pkg/database/postgres"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
)

type Repository struct {
	conn postgres.Connection
}

func NewRepository(conn postgres.Connection) *Repository {
	return &Repository{conn: conn}
}

func (repo *Repository) CreateNote(ctx context.Context, item *entity.Note) (uuid.UUID, error) {
	query, args, err := sq.Insert("notes").Columns(
		"id",
		"title",
		"payload",
		"created_at",
		"owner_id",
		"have_access",
		"draft",
		"layout_id",
	).Values(
		item.Id,
		item.Title,
		item.Payload,
		item.CreatedAt,
		item.OwnerId,
		item.HaveAccess,
		item.Draft,
		item.LayoutId,
	).PlaceholderFormat(sq.Dollar).ToSql()

	_, err = repo.conn.Exec(ctx, query, args...)
	if err != nil {
		if common.IsUniqueErr(err) {
			return uuid.Nil, apperrors.NotUnique
		}
		return uuid.Nil, err
	}
	return item.Id, err
}

func (repo *Repository) DeleteNoteById(ctx context.Context, noteId, userId uuid.UUID) error {
	query := `
		DELETE FROM notes 
		WHERE id = $1 and owner_id = $2
	`
	res, err := repo.conn.Exec(ctx, query, noteId, userId)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return apperrors.NoteNotFound
	}
	return nil
}

func (repo *Repository) DeleteNotesByLayoutId(ctx context.Context, layoutId, userId uuid.UUID) error {
	query := `
		DELETE FROM notes 
		WHERE layout_id = $1 and owner_id = $2
	`
	_, err := repo.conn.Exec(ctx, query, layoutId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) UpdateNote(ctx context.Context, newItem *entity.Note) error {
	query := `
		UPDATE notes
		SET 
		title = $1,
		payload = $2
		WHERE id = $3 and owner_id = $4
	`
	res, err := repo.conn.Exec(ctx, query, newItem.Title, newItem.Payload, newItem.Id, newItem.OwnerId)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return apperrors.NoteNotFound
	}
	return nil
}

func (repo *Repository) GetNoteCountInLayout(ctx context.Context, layoutId uuid.UUID) (int, error) {
	query := `
		select count(*) from notes where layout_id = $1
	`
	var n int
	err := repo.conn.QueryRow(ctx, query, layoutId).Scan(&n)
	return n, err
}

func (repo *Repository) UpdateDraftById(ctx context.Context, userId, noteId uuid.UUID, newDraft string) error {
	query := `
		update notes 
		set draft = $1
		where id = $2 and $3 = any(have_access)
	`
	_, err := repo.conn.Exec(ctx, query, newDraft, noteId, userId)
	return err
}

func (repo *Repository) CommitDraft(ctx context.Context, userId, noteId uuid.UUID) error {
	query := `
		update notes 
		set payload = draft, draft = ''
		where id = $1 and $2 = any(have_access)
	`
	_, err := repo.conn.Exec(ctx, query, noteId, userId)
	return err
}

// todo check access to layout
func (repo *Repository) GetNotesByLayoutId(ctx context.Context, layoutId, userId uuid.UUID, offset, limit int) ([]entity.Note, error) {
	query := `
		select n.* from notes n
		where n.layout_id = $1 and $2 = ANY(n.have_access)
		order by created_at desc, layout_id 
		offset $3
		limit $4
	`
	rows, err := repo.conn.Query(ctx, query, layoutId, userId, offset, limit)
	if err != nil {
		return nil, errors.Wrap(err, "repo.conn.Query")
	}
	defer rows.Close()

	var notes []entity.Note
	for rows.Next() {
		var item entity.Note
		err := rows.Scan(
			&item.Id,
			&item.Title,
			&item.Payload,
			&item.CreatedAt,
			&item.OwnerId,
			&item.HaveAccess,
			&item.LayoutId,
			&item.Draft,
		)
		if err != nil {
			return nil, errors.Wrap(err, "rows.Scan")
		}
		notes = append(notes, item)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows.Err")
	}
	return notes, nil
}

func (repo *Repository) GetNotesWithoutPosition(ctx context.Context, layoutId, userId uuid.UUID) ([]entity.Note, error) {
	query := `
		select n.* from notes n
		join positions p on p.note_id = n.id
		where n.layout_id = $1 and $2 = ANY(n.have_access) and x_position is null and y_position is null
	`
	rows, err := repo.conn.Query(ctx, query, layoutId, userId)
	if err != nil {
		return nil, errors.Wrap(err, "repo.conn.Query")
	}
	defer rows.Close()

	var notes []entity.Note
	for rows.Next() {
		var item entity.Note
		err := rows.Scan(
			&item.Id,
			&item.Title,
			&item.Payload,
			&item.CreatedAt,
			&item.OwnerId,
			&item.HaveAccess,
			&item.LayoutId,
			&item.Draft,
		)
		if err != nil {
			return nil, errors.Wrap(err, "rows.Scan")
		}
		notes = append(notes, item)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows.Err")
	}
	return notes, nil
}

func (repo *Repository) GetNotesWithPosition(ctx context.Context, layoutId, userId uuid.UUID) ([]entity.NoteWithPosition, error) {
	builder := sq.Select(
		"n.*",
		"p.note_id",
		"p.x_position",
		"p.y_position",
	).From("notes n").
		Join("positions p ON p.note_id = n.id").
		Where("? = ANY(n.have_access)", userId).
		Where(sq.NotEq{"p.x_position": nil}).
		Where(sq.NotEq{"p.y_position": nil}).
		PlaceholderFormat(sq.Dollar)

	if layoutId != uuid.Nil {
		builder = builder.Where(sq.Eq{"n.layout_id": layoutId})
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "builder.ToSql")
	}
	rows, err := repo.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "repo.conn.Query")
	}
	defer rows.Close()

	var notes []entity.NoteWithPosition
	for rows.Next() {
		var item entity.NoteWithPosition
		err := rows.Scan(
			&item.Id,
			&item.Title,
			&item.Payload,
			&item.CreatedAt,
			&item.OwnerId,
			&item.HaveAccess,
			&item.LayoutId,
			&item.Draft,
			&item.NoteId,
			&item.XPosition,
			&item.YPosition,
		)
		if err != nil {
			return nil, errors.Wrap(err, "rows.Scan")
		}
		notes = append(notes, item)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows.Err")
	}
	return notes, nil
}

func (repo *Repository) GetFullNotesByLayoutId(ctx context.Context, layoutId, userId uuid.UUID) ([]dto.Note, error) {
	query := `
    select 
    n.id, n.title, n.payload, n.owner_id, n.have_access, n.layout_id, n.draft,
    p.x_position, p.y_position,
    COALESCE(array((select second_note_id from links l where l.first_note_id = n.id)), '{}'),
    COALESCE(array((select first_note_id from links l where l.second_note_id = n.id)), '{}')
    from notes n
    left join positions p on p.note_id = n.id
    where n.owner_id = $1 and n.layout_id = $2
    `
	rows, err := repo.conn.Query(ctx, query, userId, layoutId)
	if err != nil {
		return nil, errors.Wrap(err, "repo.conn.Query")
	}
	defer rows.Close()

	var notes []dto.Note
	for rows.Next() {
		var item dto.Note
		var in, out pgtype.Array[uuid.UUID]

		// Временные переменные для позиции
		var xPos, yPos *float64

		err := rows.Scan(
			&item.Id,
			&item.Title,
			&item.Payload,
			&item.OwnerId,
			&item.HaveAccess,
			&item.LayoutId,
			&item.Draft,
			&xPos, // Сканируем во временную переменную
			&yPos, // Сканируем во временную переменную
			&in,
			&out,
		)
		if err != nil {
			return nil, errors.Wrap(err, "rows.Scan")
		}

		// Создаем Position только если есть координаты
		if xPos != nil && yPos != nil {
			item.Position = &dto.Position{
				XPos: *xPos,
				YPos: *yPos,
			}
		}

		if in.Valid {
			item.LinkedWithIn = in.Elements
		} else {
			item.LinkedWithIn = []uuid.UUID{}
		}

		if out.Valid {
			item.LinkedWithOut = out.Elements
		} else {
			item.LinkedWithOut = []uuid.UUID{}
		}
		notes = append(notes, item)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows.Err")
	}
	return notes, nil
}

func (repo *Repository) DeleteLinkNotes(ctx context.Context, layoutId, firstNoteId, secondNoteId uuid.UUID) error {
	query := `
		delete from links
		where layout_id = $1 and first_note_id = $2 and second_note_id = $3
	`
	_, err := repo.conn.Exec(ctx, query, layoutId, firstNoteId, secondNoteId)
	return err
}

func (repo *Repository) SearchNotes(ctx context.Context, userId uuid.UUID, search string) ([]entity.Note, error) {
	query := `
		select n.* from notes n
		where $1 = ANY(n.have_access) 
			AND (n.title ilike '%' || $2 || '%' or n.payload ilike '%' || $2 || '%')
	`
	rows, err := repo.conn.Query(ctx, query, userId, search)
	if err != nil {
		return nil, errors.Wrap(err, "repo.conn.Query")
	}
	defer rows.Close()

	var notes []entity.Note
	for rows.Next() {
		var item entity.Note
		err := rows.Scan(
			&item.Id,
			&item.Title,
			&item.Payload,
			&item.CreatedAt,
			&item.OwnerId,
			&item.HaveAccess,
			&item.LayoutId,
			&item.Draft,
		)
		if err != nil {
			return nil, errors.Wrap(err, "rows.Scan")
		}
		notes = append(notes, item)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows.Err")
	}
	return notes, nil
}
