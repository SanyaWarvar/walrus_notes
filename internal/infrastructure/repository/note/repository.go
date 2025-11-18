package note

import (
	"context"
	"wn/internal/entity"
	apperrors "wn/internal/errors"
	"wn/internal/infrastructure/repository/common"
	"wn/pkg/database/postgres"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Repository struct {
	conn postgres.Connection
}

func NewRepository(conn postgres.Connection) *Repository {
	return &Repository{conn: conn}
}

func (repo *Repository) CreateNote(ctx context.Context, item *entity.Note) (uuid.UUID, error) {
	query := `
		INSERT INTO notes VALUES
		($1, $2, $3, $4, $5, $6, &7)
		RETURNING id
	`
	var id uuid.UUID
	err := repo.conn.QueryRow(ctx, query, item.Id, item.Title, item.Payload, item.CreatedAt, item.OwnerId, item.HaveAccess, item.Draft).Scan(&id)
	if err != nil {
		if common.IsUniqueErr(err) {
			return id, apperrors.NotUnique
		}
	}
	return id, err
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

func (repo *Repository) DeleteLayoutNotes(ctx context.Context, noteId uuid.UUID) error {
	query := `
		DELETE FROM layout_note 
		WHERE note_id = $1
	`
	_, err := repo.conn.Exec(ctx, query, noteId)
	return err
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
		select count(*) from layout_note where layout_id = $1
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
		join layout_note ln on ln.note_id = n.id
		where ln.layout_id = $1 and $2 = ANY(n.have_access)
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

func (repo *Repository) DeleteLayoutNote(ctx context.Context, layoutId uuid.UUID, userId uuid.UUID) error {
	query := `
		delete from layout_note ln
		where ln.layout_id = $1 and $2 = ANY(
			select unnest(have_access) from layouts l where l.id = $1
		)
	`
	_, err := repo.conn.Exec(ctx, query, layoutId, userId)
	return err
}

func (repo *Repository) GetNotesWithoutPosition(ctx context.Context, layoutId, userId uuid.UUID) ([]entity.Note, error) {
	query := `
		select n.* from notes n
		join layout_note ln on ln.note_id = n.id
		where ln.layout_id = $1 and $2 = ANY(n.have_access) and x_position is null and y_position is null
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
	query := `
		select n.*, ln.note_id, ln.layout_id, ln.x_position, ln.y_position from notes n
		join layout_note ln on ln.note_id = n.id
		where ln.layout_id = $1 and $2 = ANY(n.have_access) and x_position is not null and y_position is not null
	`
	rows, err := repo.conn.Query(ctx, query, layoutId, userId)
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
			&item.NoteId,
			&item.LayoutId,
			&item.XPosition,
			&item.YPosition,
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

func (repo *Repository) UpdateNotePosition(ctx context.Context, layoutId, noteId uuid.UUID, xPos, yPos *float64) error {
	query := `
		update layout_note ln
		set x_position = $1, y_position = $2
		where ln.layout_id = $3 and ln.note_id = $4
	`
	_, err := repo.conn.Exec(ctx, query, xPos, yPos, layoutId, noteId)
	return err
}

func (repo *Repository) LinkNotes(ctx context.Context, layoutId, firstNoteId, secondNoteId uuid.UUID) error {
	query := `
		insert into links
		values ($1, $2, $3)
	`
	_, err := repo.conn.Exec(ctx, query, layoutId, firstNoteId, secondNoteId)
	return err
}

func (repo *Repository) DeleteLinkNotes(ctx context.Context, layoutId, firstNoteId, secondNoteId uuid.UUID) error {
	query := `
		delete from links
		where layout_id = $1 and first_note_id = $2 and second_note_id = $3
	`
	_, err := repo.conn.Exec(ctx, query, layoutId, firstNoteId, secondNoteId)
	return err
}

func (repo *Repository) DeleteLinksFromLayout(ctx context.Context, layoutId uuid.UUID) error {
	query := `
		delete from links
		where layout_id = $1
	`
	_, err := repo.conn.Exec(ctx, query, layoutId)
	return err
}

func (repo *Repository) DeleteLinksWithNote(ctx context.Context, noteId uuid.UUID) error {
	query := `
		delete from links
		where first_note_id = $1 or second_note_id = $1
	`
	_, err := repo.conn.Exec(ctx, query, noteId)
	return err
}

func (repo *Repository) GetAllLinks(ctx context.Context, layoutId uuid.UUID, noteIds []uuid.UUID) ([]entity.Link, error) {
	query := `
		select layout_id, first_note_id, second_note_id
		from links
		where layout_id = $1 and (first_note_id = ANY($2) or second_note_id = ANY($2))
	`
	rows, err := repo.conn.Query(ctx, query, layoutId, noteIds)
	if err != nil {
		return nil, errors.Wrap(err, "repo.conn.Query")
	}
	defer rows.Close()

	var links []entity.Link
	for rows.Next() {
		var link entity.Link
		err := rows.Scan(
			&link.LayoutId,
			&link.FirstNoteId,
			&link.SecondNoteId,
		)
		if err != nil {
			return nil, errors.Wrap(err, "rows.Scan")
		}
		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows.Err")
	}
	return links, nil
}

func (repo *Repository) SearchNotes(ctx context.Context, userId uuid.UUID, search string) ([]entity.Note, error) {
	query := `
		select n.* from notes n
		join layout_note ln on ln.note_id = n.id
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
