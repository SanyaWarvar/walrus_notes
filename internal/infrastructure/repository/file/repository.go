package file

import (
	"context"
	"fmt"
	"wn/pkg/database/postgres"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

type Repository struct {
	conn postgres.Connection
}

func NewRepository(conn postgres.Connection) *Repository {
	return &Repository{conn: conn}
}

func (repo *Repository) CreateFile(ctx context.Context, filename string, encodedFile string) error {
	query, args, err := squirrel.Insert("files").Values(filename, encodedFile).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return errors.Wrap(err, "squirrel.ToSql")
	}

	_, err = repo.conn.Exec(ctx, query, args...)
	return err
}

func (repo *Repository) GetAllFiles(ctx context.Context) ([]StaticFile, error) {

	query, args, err := squirrel.Select("file_name, file_data").From("files").PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "squirrel.ToSql")
	}

	rows, err := repo.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "repo.conn.Query")
	}
	defer rows.Close()
	var files []StaticFile
	for rows.Next() {
		var f StaticFile
		if err := rows.Scan(&f.Filename, &f.FileAsString); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		files = append(files, f)
	}
	return files, nil
}
