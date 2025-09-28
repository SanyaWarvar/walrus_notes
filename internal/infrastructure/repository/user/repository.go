package user

import (
	"context"
	"database/sql"
	apperrors "wn/internal/errors"
	"wn/internal/infrastructure/repository/common"
	"wn/pkg/database/postgres"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/Masterminds/squirrel"
)

type Repository struct {
	conn postgres.Connection
}

func NewRepository(conn postgres.Connection) *Repository {
	return &Repository{conn: conn}
}

func (repo *Repository) CreateUser(ctx context.Context, item *User) error {
	query := `
		INSERT INTO users VALUES
		($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := repo.conn.Exec(ctx, query, item.Id, item.Username, item.Email, item.Password, item.Role, item.ImgUrl, item.ConfirmedEmail, item.CreatedAt)
	if err != nil {
		if common.IsUniqueErr(err) {
			return apperrors.NotUnique
		}
	}
	return err
}

func (repo *Repository) UpdateUser(ctx context.Context, userId uuid.UUID, updateParams *UserUpdateParams) error {

	builder := squirrel.Update("users").
		Where(squirrel.Eq{"id": userId}).
		PlaceholderFormat(squirrel.Dollar)

	if updateParams.Username != nil {
		builder = builder.Set("username", *updateParams.Username)
	}
	if updateParams.Email != nil {
		builder = builder.Set("email", *updateParams.Email)
	}
	if updateParams.Password != nil {
		builder = builder.Set("password", *updateParams.Password)
	}
	if updateParams.ImgUrl != nil {
		builder = builder.Set("img_url", *updateParams.ImgUrl)
	}
	if updateParams.ConfirmedEmail != nil {
		builder = builder.Set("confirmed_email", *updateParams.ConfirmedEmail)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "builder.ToSql")
	}

	_, err = repo.conn.Exec(ctx, query, args...)
	if err != nil {
		if common.IsUniqueErr(err) {
			return apperrors.NotUnique
		}
		return errors.Wrap(err, "repo.conn.Exec")
	}

	return nil
}

func (repo *Repository) GetUser(ctx context.Context, filter UserFilter) (*User, bool, error) {
	var output User
	builder := squirrel.Select("u.*").From("users u")

	if filter.Id != nil {
		builder = builder.Where(squirrel.Eq{"id": filter.Id})
	}

	if filter.Email != nil {
		builder = builder.Where(squirrel.Eq{"email": filter.Email})
	}

	if filter.Limit > 0 {
		builder.Limit(filter.Limit)
	}

	query, args, err := builder.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, false, errors.Wrap(err, "squirrel.ToSql")
	}
	err = repo.conn.QueryRow(ctx, query, args...).
		Scan(&output.Id, &output.Username, &output.Email, &output.Password, &output.Role, &output.ImgUrl, &output.ConfirmedEmail, &output.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &output, false, nil
		}
		if common.IsUniqueErr(err) {
			return &output, false, apperrors.NotUnique
		}
		return &output, false, errors.Wrap(err, "repo.conn.QueryRow.Scan")
	}

	return &output, true, nil
}
