package tokens

import (
	"context"
	"database/sql"
	"wn/pkg/database/postgres"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Repository struct {
	conn postgres.Connection
}

func NewRepository(conn postgres.Connection) *Repository {
	return &Repository{conn: conn}
}

// Create создает новый refresh token
func (repo *Repository) Create(ctx context.Context, token *RefreshToken) error {
	query, args, err := squirrel.Insert("refresh_tokens").
		Columns("id", "user_id", "access_id", "exp_at").
		Values(token.Id, token.UserId, token.AccessId, token.ExpAt).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "squirrel.ToSql")
	}

	_, err = repo.conn.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "repo.conn.Exec")
	}

	return nil
}

// GetByID возвращает refresh token по ID
func (repo *Repository) GetByID(ctx context.Context, id uuid.UUID) (*RefreshToken, bool, error) {
	query, args, err := squirrel.Select("id", "user_id", "access_id", "exp_at").
		From("refresh_tokens").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, false, errors.Wrap(err, "squirrel.ToSql")
	}

	var token RefreshToken
	err = repo.conn.QueryRow(ctx, query, args...).Scan(&token.Id, &token.UserId, &token.AccessId, &token.ExpAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &token, false, nil
		}
		return nil, false, errors.Wrap(err, "repo.conn.QueryRow.Scan")
	}

	return &token, true, nil
}

// Delete удаляет refresh token по ID
func (repo *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := squirrel.Delete("refresh_tokens").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "squirrel.ToSql")
	}

	_, err = repo.conn.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "repo.conn.Exec")
	}

	return nil
}

// todo add
func (repo *Repository) DeleteExpired(ctx context.Context, cutoffTime time.Time) error {
	query, args, err := squirrel.Delete("refresh_tokens").
		Where(squirrel.Lt{"exp_at": cutoffTime}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "squirrel.ToSql")
	}

	_, err = repo.conn.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "repo.conn.Exec")
	}

	return nil
}
