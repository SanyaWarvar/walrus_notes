package note

import (
	"context"
	"wn/internal/entity"
	apperrors "wn/internal/errors"
	"wn/internal/infrastructure/repository/common"
	"wn/pkg/database/postgres"

	"github.com/google/uuid"
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
		($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	var id uuid.UUID
	err := repo.conn.QueryRow(ctx, query, item.Id, item.Title, item.Payload, item.CreatedAt, item.OwnerId, item.HaveAccess).Scan(&id)
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
