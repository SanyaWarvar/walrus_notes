package links

import (
	"context"
	"wn/internal/entity"
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

func (repo *Repository) LinkNotes(ctx context.Context, firstNoteId, secondNoteId uuid.UUID) error {
	query := `
		insert into links
		values ($1, $2)
	`
	_, err := repo.conn.Exec(ctx, query, firstNoteId, secondNoteId)
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

func (repo *Repository) DeleteLinksByLayoutId(ctx context.Context, layoutId uuid.UUID) error {
	query := `
		delete from links
		where first_note_id = any(
			select id from notes n where n.layout_id = $1
		) or second_note_id = any(
			select id from notes n where n.layout_id = $1
		)
	`
	_, err := repo.conn.Exec(ctx, query, layoutId)
	return err
}

func (repo *Repository) GetAllLinks(ctx context.Context, noteIds []uuid.UUID) ([]entity.Link, error) {
	query := `
		select first_note_id, second_note_id
		from links
		where first_note_id = ANY($1) or second_note_id = ANY($1)
	`
	rows, err := repo.conn.Query(ctx, query, noteIds)
	if err != nil {
		return nil, errors.Wrap(err, "repo.conn.Query")
	}
	defer rows.Close()

	var links []entity.Link
	for rows.Next() {
		var link entity.Link
		err := rows.Scan(
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

func (repo *Repository) DeleteLink(ctx context.Context, noteId1, noteId2 uuid.UUID) error {
	query := `
		delete from links
		where first_note_id = $1 and second_note_id = $2
	`
	_, err := repo.conn.Exec(ctx, query, noteId1, noteId2)
	return err
}
