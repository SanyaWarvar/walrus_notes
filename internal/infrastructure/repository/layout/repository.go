package layout

import (
	"context"
	"wn/internal/entity"
	apperrors "wn/internal/errors"
	"wn/internal/infrastructure/repository/common"
	"wn/pkg/database/postgres"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
		($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	var id uuid.UUID
	err := repo.conn.QueryRow(ctx, query, item.Id, item.Title, item.OwnerId, item.HaveAccess, item.IsMain, item.Color).Scan(&id)
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
		return apperrors.LayoutNotFound
	}
	return nil
}

func (repo *Repository) GetAvailableLayouts(ctx context.Context, userId uuid.UUID) ([]entity.Layout, error) {
	query := `
		select l.* 
		from layouts l
		WHERE $1 = l.owner_id or $1 = any(select to_user_id from permissions p where p.target_id = l.id and p.to_user_id = $1)
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
			&layout.Color,
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

func (repo *Repository) UpdateLayout(ctx context.Context, userId, layoutId uuid.UUID, color, title string) (int, error) {
	builder := squirrel.Update("layouts ").
		Where(squirrel.Eq{"id": layoutId}).
		Where("? = ANY(have_access)", userId).
		PlaceholderFormat(squirrel.Dollar)

	if color != "" {
		builder = builder.Set("color", color)
	}
	if title != "" {
		builder = builder.Set("title", title)
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "builder.ToSql")
	}
	res, err := repo.conn.Exec(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "repo.conn.Exec")
	}
	return int(res.RowsAffected()), nil
}

func (repo *Repository) GetByOwnerId(ctx context.Context, ownerId, layoutId uuid.UUID) (*entity.Layout, error) {
	sql, args, err := sq.
		Select("l.*").
		From("layouts l").
		Where(sq.Eq{"owner_id": ownerId}).
		Where(sq.Eq{"id": layoutId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "toSql")
	}

	var item entity.Layout
	err = repo.conn.QueryRow(ctx, sql, args...).Scan(
		&item.Id,
		&item.Title,
		&item.OwnerId,
		&item.HaveAccess,
		&item.IsMain,
		&item.Color,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.RecordNotFound
		}
		return nil, errors.Wrap(err, "scan")
	}

	return &item, nil
}

func (repo *Repository) GetById(ctx context.Context, layoutId uuid.UUID) (*entity.Layout, error) {
	sql, args, err := sq.
		Select("l.*").
		From("layouts l").
		Where(sq.Eq{"id": layoutId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "toSql")
	}

	var item entity.Layout
	err = repo.conn.QueryRow(ctx, sql, args...).Scan(
		&item.Id,
		&item.Title,
		&item.OwnerId,
		&item.HaveAccess,
		&item.IsMain,
		&item.Color,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.RecordNotFound
		}
		return nil, errors.Wrap(err, "scan")
	}

	return &item, nil
}
