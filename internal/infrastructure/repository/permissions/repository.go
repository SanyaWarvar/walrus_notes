package permissions

import (
	"context"
	"wn/internal/domain/dto"
	"wn/internal/entity"
	apperrors "wn/internal/errors"
	"wn/pkg/database/postgres"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Repository struct {
	conn postgres.Connection
}

func NewRepository(conn postgres.Connection) *Repository {
	return &Repository{conn: conn}
}

func (repo *Repository) CreatePermissions(ctx context.Context, item *entity.Permission) error {
	query := sq.
		Insert("permissions").
		Columns(
			"id",
			"to_user_id",
			"from_user_id",
			"target_id",
			"kind",
			"can_read",
			"can_write",
			"can_edit",
			"created_at",
		).Values(
		item.Id,
		item.ToUserId,
		item.FromUserId,
		item.TargetId,
		item.Kind,
		item.CanRead,
		item.CanWrite,
		item.CanEdit,
		item.CreatedAt,
	).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "ToSql")
	}

	_, err = repo.conn.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "Exec")
	}
	return nil
}

func (repo *Repository) UpdatePermissions(ctx context.Context, item *entity.Permission) error {
	query := sq.
		Update("permissions").
		Set("can_read", item.CanRead).
		Set("can_write", item.CanWrite).
		Set("can_edit", item.CanEdit).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "ToSql")
	}

	_, err = repo.conn.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "Exec")
	}
	return nil
}

func (repo *Repository) DeletePermissions(ctx context.Context, permissionsIds ...uuid.UUID) error {
	query := sq.
		Delete("permissions").
		Where(sq.Eq{"id": permissionsIds}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return errors.Wrap(err, "ToSql")
	}

	_, err = repo.conn.Exec(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, "Exec")
	}
	return nil
}

func (repo *Repository) GetPermissions(ctx context.Context, filter *dto.GetPermissionsFilter) ([]entity.Permission, error) {
	query := sq.
		Select(
			"id",
			"to_user_id",
			"from_user_id",
			"target_id",
			"kind",
			"can_read",
			"can_write",
			"can_edit",
			"created_at",
		).
		From("permissions p").
		PlaceholderFormat(sq.Dollar)

	if filter.Id != nil {
		query = query.Where(sq.Eq{"p.id": filter.Id})
	}

	if filter.Kind != nil {
		query = query.Where(sq.Eq{"p.kind": filter.Kind})
	}

	if filter.FromUserId != nil {
		query = query.Where(sq.Eq{"p.from_user_id": filter.FromUserId})
	}

	if filter.ToUserId != nil {
		query = query.Where(sq.Eq{"p.to_user_id": filter.ToUserId})
	}

	if filter.TargetId != nil {
		query = query.Where(sq.Eq{"p.target_id": filter.TargetId})
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "ToSql")
	}

	rows, err := repo.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "repo.conn.Query")
	}
	defer rows.Close()

	var items []entity.Permission
	for rows.Next() {
		var item entity.Permission
		err := rows.Scan(
			&item.Id,
			&item.ToUserId,
			&item.FromUserId,
			&item.TargetId,
			&item.Kind,
			&item.CanRead,
			&item.CanWrite,
			&item.CanEdit,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "rows.Scan")
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows.Err")
	}

	return items, nil
}

func (repo *Repository) GetPermission(ctx context.Context, filter *dto.GetPermissionsFilter) (*entity.Permission, error) {
	filter.Limit = 1

	items, err := repo.GetPermissions(ctx, filter)
	if err != nil {
		return nil, err
	}

	if len(items) > 0 {
		return &items[0], nil
	}

	return nil, apperrors.RecordNotFound
}
