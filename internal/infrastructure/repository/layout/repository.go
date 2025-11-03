package layout

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

func (repo *Repository) CreateLayout(ctx context.Context, item *entity.Layout) (uuid.UUID, error) {
	query := `
		INSERT INTO layouts VALUES
		($1, $2, $3, $4, $5)
		RETURNING id
	`
	var id uuid.UUID
	err := repo.conn.QueryRow(ctx, query, item.Id, item.Title, item.OwnerId, item.HaveAccess, item.IsMain).Scan(&id)
	if err != nil {
		if common.IsUniqueErr(err) {
			return id, apperrors.NotUnique
		}
	}
	return id, err
}

func (repo *Repository) DeleteLayoutById(ctx context.Context, layoutId, userId uuid.UUID) error {
	query := `
		DELETE FROM layouts 
		WHERE id = $1 and owner_id = $2
	`
	res, err := repo.conn.Exec(ctx, query, layoutId, userId)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return apperrors.NoteNotFound
	}
	return nil
}

func (repo *Repository) GetAvailableLayouts(ctx context.Context, userId uuid.UUID) ([]entity.Layout, error) {
	query := `
		select * 
		from layouts l
		WHERE $1 = ANY(l.have_access)
	`
	rows, err := repo.conn.Query(ctx, query, userId)
	if err != nil {
		return nil, errors.Wrap(err, "repo.conn.Query")
	}
	defer rows.Close()

	var layouts []entity.Layout
	for rows.Next() {
		var layout entity.Layout
		err := rows.Scan(
			&layout.Id,
			&layout.Title,
			&layout.OwnerId,
			&layout.HaveAccess,
			&layout.IsMain,
		)
		if err != nil {
			return nil, errors.Wrap(err, "rows.Scan")
		}
		layouts = append(layouts, layout)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows.Err")
	}

	return layouts, nil
}

func (repo *Repository) LinkNoteWithLayout(ctx context.Context, noteId, layoutId uuid.UUID) error {
	query := `
		INSERT INTO layout_note VALUES
		($1, $2, $3, $4)
	`
	_, err := repo.conn.Exec(ctx, query, noteId, layoutId, nil, nil)
	return err
}
