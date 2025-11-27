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

func (repo *Repository) UpdateNotePosition(ctx context.Context, noteId uuid.UUID, xPos, yPos *float64) error {
	query := `
		update positions p
		set x_position = $1, y_position = $2
		where p.note_id = $3
	`
	_, err := repo.conn.Exec(ctx, query, xPos, yPos, noteId)
	return err
}
