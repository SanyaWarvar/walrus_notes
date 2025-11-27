package positions

import (
	"context"
	"wn/pkg/database/postgres"

	"github.com/google/uuid"
)

type Repository struct {
	conn postgres.Connection
}

func NewRepository(conn postgres.Connection) *Repository {
	return &Repository{conn: conn}
}

func (repo *Repository) CreateNotePosition(ctx context.Context, noteId uuid.UUID, xPos, yPos *float64) error {
	query := `
		insert into positions p (note_id, x_position, y_position)
		values ($1, $2, $3)
	`
	_, err := repo.conn.Exec(ctx, query, noteId, xPos, yPos)
	return err
}

func (repo *Repository) UpdateNotePosition(ctx context.Context, noteId uuid.UUID, xPos, yPos *float64) error {
	query := `
		update positions p
		set x_position = $1, y_position = $2
		where p.note_id = $3
	`
	_, err := repo.conn.Exec(ctx, query, xPos, yPos, noteId)
	return err
}
